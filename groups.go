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

func (c *Client) GetGroups(ctx context.Context, params GetGroupsDefinition) (GetGroupsResponse, error) {
	raw := GetGroupsRequest{
		Method: "get",
		Params: params,
	}

	req, err := c.createGetGroupsRequest(ctx, raw)
	if err != nil {
		return GetGroupsResponse{}, fmt.Errorf("createGetReportRequest: %w", err)
	}

	time.Sleep(time.Duration(c.statisticsLimit.retryInterval) * time.Second)

	resp, err := c.Tr.Do(req)
	if err != nil {
		return GetGroupsResponse{}, fmt.Errorf("do request: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return GetGroupsResponse{}, err
	}

	err = statusCodeHandler(resp)
	if err != nil {
		return GetGroupsResponse{}, fmt.Errorf("statusCodeHandler: %w", err)
	}

	var errResp APIError

	err = json.Unmarshal(body, &errResp)
	if err != nil {
		return GetGroupsResponse{}, err
	}

	if errResp.Err.ErrorCode != 0 {
		return GetGroupsResponse{}, errResp
	}

	var response GetGroupsResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return GetGroupsResponse{}, err
	}

	return response, nil
}

func (c *Client) createGetGroupsRequest(ctx context.Context, params GetGroupsRequest) (*http.Request, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.direct.yandex.com/json/v5/adgroups", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	c.buildHeader(req)

	return req, nil
}

type GetGroupsRequest struct {
	Method string              `json:"method"`
	Params GetGroupsDefinition `json:"params"`
}

type GetGroupsDefinition struct {
	Selection                    GroupsSelectionCriteria `json:"SelectionCriteria,omitempty"`
	FieldNames                   []string                `json:"FieldNames"`
	DynamicTextAdGroupFieldNames []string                `json:"DynamicTextAdGroupFieldNames,omitempty"`
}

type GroupsSelectionCriteria struct {
	CampaignIds     []int    `json:"CampaignIds,omitempty"`
	IDs             []int    `json:"Ids,omitempty"`
	Statuses        []string `json:"Statuses,omitempty"`
	ServingStatuses []string `json:"ServingStatuses,omitempty"`
}

type GetGroupsResponse struct {
	Result struct {
		AdGroups []AdGroupItem `json:"AdGroups"`
	} `json:"result"`
}

type AdGroupItem struct {
	ID            int    `json:"Id"`
	Name          string `json:"Name"`
	CampaignID    int64  `json:"CampaignId"`
	Status        string `json:"Status"`
	ServingStatus string `json:"ServingStatus"`
	Type          string `json:"Type"`
}
