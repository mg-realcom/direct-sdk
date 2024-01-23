package direct_sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	anti_fraud "github.com/mg-realcom/direct-sdk/anti-fraud"
	"github.com/mg-realcom/direct-sdk/common"
	"github.com/mg-realcom/direct-sdk/statistics"
	"github.com/rs/zerolog"
	"io"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	Tr              *http.Client
	Login           string
	Token           *string
	App             *App
	host            environment
	statisticsLimit statisticsLimits
	logger          *zerolog.Logger
}

type App struct {
	ID     string
	Secret string
}

type statisticsLimits struct {
	retryInterval  int32
	reportsInQueue int8
}

type environment string

const (
	LIVE    environment = "api.direct.yandex.com"
	SANDBOX environment = "api-sandbox.direct.yandex.com"
)

func NewClient(tr *http.Client, login string, token *string, app *App, sandbox bool, logger *zerolog.Logger) *Client {
	if sandbox {
		return &Client{
			Login: login,
			Token: token,
			host:  SANDBOX,
		}
	}

	return &Client{
		Tr:    tr,
		Login: login,
		Token: token,
		App:   app,
		host:  LIVE,
		statisticsLimit: statisticsLimits{
			retryInterval:  0,
			reportsInQueue: 0,
		},
		logger: logger,
	}
}

func (c *Client) buildHeader(req *http.Request) {
	req.Header.Add("Authorization", "Bearer "+*c.Token)
	req.Header.Add("Client-Login", c.Login)
	req.Header.Add("Accept-Language", "ru")
	req.Header.Add("skipReportHeader", "true")
	req.Header.Add("skipReportSummary", "true")
}

type Payload struct {
	Method string `json:"method"`
	Params struct {
		Ads []struct {
			AdGroupID int `json:"AdGroupId"`
			TextAd    struct {
				Text   string `json:"Text"`
				Title  string `json:"Title"`
				Href   string `json:"Href,omitempty"`
				Mobile string `json:"Mobile"`
			} `json:"TextAd"`
		} `json:"Ads"`
	} `json:"params"`
}

func (c *Client) GetCampaigns(ctx context.Context, dateRange statistics.DateRange) ([]string, error) {
	rand.Seed(time.Now().UnixNano())
	randomInt := rand.Intn(99999)
	repName := fmt.Sprintf("campaigns_%s_%s - %s_%d", c.Login, dateRange.From, dateRange.To, randomInt)
	fmt.Println(repName)
	params := statistics.ReportDefinition{
		Selection: &statistics.SelectionCriteria{
			DateFrom: dateRange.From,
			DateTo:   dateRange.To,
			Filter:   nil,
		},
		DateRangeType: statistics.DateRangeCustomDate,
		ReportType:    statistics.CustomReport,
		FieldNames:    []string{"CampaignType"},
		ReportName:    repName,
		Format:        common.FormatTSV,
		IncludeVAT:    common.NO,
	}
	reqContent := Request{Params: params}
	body, err := json.Marshal(reqContent)
	if err != nil {
		return nil, err
	}
	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.direct.yandex.com/json/v5/reports", bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}
		c.buildHeader(req)
		c.waitInfo(params.ReportName)
		time.Sleep(time.Duration(c.statisticsLimit.retryInterval) * time.Second)
		resp, err := c.Tr.Do(req)
		defer resp.Body.Close()

		if err != nil {
			return nil, fmt.Errorf("do request: %w", err)
		}
		respDump, _ := httputil.DumpResponse(resp, true)
		switch resp.StatusCode {
		case http.StatusOK:
			responseBody, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			campaignsSlice := strings.Fields(string(responseBody))
			if len(campaignsSlice) < 2 {
				fmt.Println(campaignsSlice)
				return nil, errors.New("no campaigns")
			}
			return campaignsSlice[1:], nil
		case http.StatusCreated, http.StatusAccepted:
			err := c.waitInit(resp)
			if err != nil {
				return nil, fmt.Errorf("waitInit: %w", err)
			}
		case http.StatusInternalServerError:
			c.logger.Info().Msg(fmt.Sprintf("RESPONSE:\n%s", respDump))
			return nil, errors.New("internal server error")
		case http.StatusBadRequest:
			data, err := c.badRequestPrepare(resp)
			if err != nil {
				return nil, fmt.Errorf("cannot prepare bad request: %w", err)
			}
			return nil, fmt.Errorf("ошибка получения кампаний %s", data.Error.ErrorDetail)
		default:
			return nil, fmt.Errorf("cтатус код сервера при получении кампаний %v", resp.StatusCode)
		}
	}
}

func (c *Client) GetReport(ctx context.Context, prefixTitleRequest, dir string, typeReport statistics.ReportType, fields []string, filter []statistics.Filter, dateRange statistics.DateRange) ([]string, error) {
	t := time.Now().Format("2006-01-02")
	var reportName string
	var dtRangeType statistics.DateRangeType

	reportName = fmt.Sprintf("%s_%s_%s_%s_%s", c.Login, prefixTitleRequest, dateRange.From, dateRange.To, t)
	dtRangeType = statistics.DateRangeCustomDate
	params := statistics.ReportDefinition{
		Selection: &statistics.SelectionCriteria{
			DateFrom: dateRange.From,
			DateTo:   dateRange.To,
			Filter:   filter,
		},
		FieldNames:    fields,
		Page:          &common.Page{Limit: 50_000, Offset: 0},
		ReportName:    reportName,
		ReportType:    typeReport,
		DateRangeType: dtRangeType,
		Format:        common.FormatTSV,
		IncludeVAT:    common.NO,
	}

	fileNames, err := c.GetFiles(ctx, dir, params)
	if err != nil {
		return nil, fmt.Errorf("GetFiles for %s – %w", c.Login, err)
	}
	return fileNames, nil
}

func (c *Client) GetFiles(ctx context.Context, dir string, params statistics.ReportDefinition) ([]string, error) {
	var result []string
	part := 1
	reportName := params.ReportName
	params.ReportName += fmt.Sprintf("_part_%d", part)
	fieldsSize := 0
	for _, field := range params.FieldNames {
		fieldsSize += len([]byte(field))
	}
	fieldsSize += len(params.FieldNames)
	for {
		req, err := c.createGetReportRequest(ctx, params)
		if err != nil {
			return result, fmt.Errorf("createGetReportRequest: %w", err)
		}
		reqDump, _ := httputil.DumpRequestOut(req, true)
		c.waitInfo(params.ReportName)
		time.Sleep(time.Duration(c.statisticsLimit.retryInterval) * time.Second)
		resp, err := c.Tr.Do(req)
		if err != nil {
			return result, fmt.Errorf("do request: %w", err)
		}
		respDump, _ := httputil.DumpResponse(resp, true)

		switch resp.StatusCode {
		case http.StatusOK:
			file, fileSize, err := createTSVFile(dir, params.ReportName, resp)
			if err != nil {
				return result, fmt.Errorf("createTSVFile: %w", err)
			}
			if fieldsSize < fileSize {
				result = append(result, file)
				params.Page.Offset += params.Page.Limit
				part++
				params.ReportName = reportName + fmt.Sprintf("_part_%d", part)
				continue
			} else {
				_ = os.Remove(file)
				return result, nil
			}
		case http.StatusCreated, http.StatusAccepted:
			err := c.waitInit(resp)
			if err != nil {
				return result, fmt.Errorf("waitInit: %w", err)
			}
		case http.StatusInternalServerError:
			c.logger.Info().Msg(fmt.Sprintf("REQUEST:\n%s", reqDump))
			c.logger.Info().Msg(fmt.Sprintf("RESPONSE:\n%s", respDump))
			return result, errors.New("internal server error")
		case http.StatusBadRequest:
			data, err := c.badRequestPrepare(resp)
			if err != nil {
				return result, fmt.Errorf("cannot prepare bad request: %w", err)
			}
			return result, fmt.Errorf("ошибка отчета %s", data.Error.ErrorDetail)
		default:
			return result, fmt.Errorf("cтатус код сервера при получении отчета %v", resp.StatusCode)
		}
	}
}

type Request struct {
	Params statistics.ReportDefinition `json:"params"`
}

type Response struct {
	Error struct {
		ErrorDetail string `json:"error_detail"`
		RequestID   string `json:"request_id"`
		ErrorCode   string `json:"error_code"`
		ErrorString string `json:"error_string"`
	} `json:"error"`
}

func (c *Client) createGetReportRequest(ctx context.Context, params statistics.ReportDefinition) (*http.Request, error) {
	reqContent := Request{Params: params}
	body, err := json.Marshal(reqContent)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.direct.yandex.com/json/v5/reports", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	c.buildHeader(req)

	return req, nil
}

func (c *Client) badRequestPrepare(resp *http.Response) (Response, error) {
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("cant read response body: %w", err)
	}

	var data Response

	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		return Response{}, fmt.Errorf("cant unmarshal response body: %w", err)
	}

	return data, nil
}

func (c *Client) waitInit(resp *http.Response) error {
	if resp == nil {
		return fmt.Errorf("response is nil")
	}

	retryIn, err := strconv.Atoi(resp.Header.Get("retryIn"))
	if err != nil {
		return fmt.Errorf("retryIn: %w", err)
	}

	c.statisticsLimit.retryInterval = int32(retryIn)

	reportsInQueue, err := strconv.Atoi(resp.Header.Get("reportsInQueue"))
	if err != nil {
		return fmt.Errorf("reportsInQueue: %w", err)
	}

	c.statisticsLimit.reportsInQueue = int8(reportsInQueue)

	return nil
}

func (c *Client) waitInfo(reportName string) {
	if c.statisticsLimit.retryInterval > 1 {
		c.logger.Info().Msg(fmt.Sprintf("Повтор запроса на отчет %s через %v\n", reportName, c.statisticsLimit.retryInterval))
	}

	if c.statisticsLimit.reportsInQueue > 1 {
		c.logger.Info().Msg(fmt.Sprintf("Количество отчетов в очереди %v\n", c.statisticsLimit.reportsInQueue))
	}
}

func (c *Client) GetYclidStat(ctx context.Context, yclids []string) ([]anti_fraud.Row, error) {

	reqs := []anti_fraud.Requests{}
	for _, v := range yclids {
		req := anti_fraud.Requests{
			Yclid: v,
		}
		reqs = append(reqs, req)
	}
	reqContent := anti_fraud.Req{Method: "get", Params: anti_fraud.Params{FieldNames: []string{"Yclid", "Phone", "Email", "Score"},
		SelectionCriteria: anti_fraud.SelectionCriteria{Requests: reqs}}}
	body, err := json.Marshal(reqContent)
	if err != nil {
		fmt.Println(err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.direct.yandex.com/json/v5/conversionscores", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", "Bearer "+*c.Token)
	req.Header.Add("Client-Login", c.Login)
	resp, err := c.Tr.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	t := time.Now().Format(time.DateTime)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}

	var res anti_fraud.Response
	reader, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(reader, &res)

	for i := 0; i < len(res.Result.ConversionScores); i++ {
		res.Result.ConversionScores[i].Date = t
		res.Result.ConversionScores[i].Login = c.Login
	}
	return res.Result.ConversionScores, nil
}

func createTSVFile(dir string, filename string, resp *http.Response) (string, int, error) {

	if resp == nil {
		return "", 0, fmt.Errorf("response is nil")
	}
	defer resp.Body.Close()

	f, err := os.CreateTemp(dir, fmt.Sprintf("%s_*.tsv", filename))
	if err != nil {
		return "", 0, fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return "", 0, err
	}
	stat, _ := f.Stat()
	size := stat.Size()
	return f.Name(), int(size), nil
}
