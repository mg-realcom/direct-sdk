package direct_sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func (c *Client) GetCampaigns(ctx context.Context, params GetCampaignsDefinition) (GetCampaignsResponse, error) {
	raw := GetCampaignsRequest{
		Method: "get",
		Params: params,
	}

	req, err := c.createGetCampaignsRequest(ctx, raw)
	if err != nil {
		return GetCampaignsResponse{}, fmt.Errorf("createGetReportRequest: %w", err)
	}

	time.Sleep(time.Duration(c.statisticsLimit.retryInterval) * time.Second)

	resp, err := c.Tr.Do(req)
	if err != nil {
		return GetCampaignsResponse{}, fmt.Errorf("do request: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return GetCampaignsResponse{}, err
	}

	err = statusCodeHandler(resp)
	if err != nil {
		return GetCampaignsResponse{}, fmt.Errorf("statusCodeHandler: %w", err)
	}

	var errResp APIError

	err = json.Unmarshal(body, &errResp)
	if err != nil {
		return GetCampaignsResponse{}, err
	}

	if errResp.Err.ErrorCode != 0 {
		return GetCampaignsResponse{}, errResp
	}

	var response GetCampaignsResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return GetCampaignsResponse{}, err
	}

	return response, nil
}

func (c *Client) createGetCampaignsRequest(ctx context.Context, params GetCampaignsRequest) (*http.Request, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.direct.yandex.com/json/v5/campaigns", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	c.buildHeader(req)

	return req, nil
}

type GetCampaignsRequest struct {
	Method string                 `json:"method"`
	Params GetCampaignsDefinition `json:"params"`
}

type GetCampaignsDefinition struct {
	Selection              CampaignsSelectionCriteria `json:"SelectionCriteria,omitempty"`
	FieldNames             []string                   `json:"FieldNames"`
	TextCampaignFieldNames *[]string                  `json:"TextCampaignFieldNames,omitempty"`
}

type CampaignsSelectionCriteria struct {
	IDs             *[]int    `json:"Ids,omitempty"`
	Types           *[]string `json:"Types,omitempty"`
	States          *[]string `json:"States,omitempty"`
	Statuses        *[]string `json:"Statuses,omitempty"`
	StatusesPayment *[]string `json:"StatusesPayment,omitempty"`
}

type GetCampaignsResponse struct {
	Result struct {
		Campaigns []CampaignItem `json:"Campaigns"`
	} `json:"result"`
}

type CampaignItem struct {
	ID                  int    `json:"Id"`
	Name                string `json:"Name"`
	StartDate           string `json:"StartDate"`
	Type                string `json:"Type"`
	Status              string `json:"Status"`
	State               string `json:"State"`
	StatusPayment       string `json:"StatusPayment"`
	StatusClarification string `json:"StatusClarification"`
	SourceId            int    `json:"SourceId"`
}
