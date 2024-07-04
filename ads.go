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

func (c *Client) GetAds(ctx context.Context, params GetAdsDefinition) (GetAdsResponse, error) {
	raw := GetAdsRequest{
		Method: "get",
		Params: params,
	}

	req, err := c.createGetAdsRequest(ctx, raw)
	if err != nil {
		return GetAdsResponse{}, fmt.Errorf("createGetReportRequest: %w", err)
	}

	time.Sleep(time.Duration(c.statisticsLimit.retryInterval) * time.Second)

	resp, err := c.Tr.Do(req)
	if err != nil {
		return GetAdsResponse{}, fmt.Errorf("do request: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return GetAdsResponse{}, err
	}

	err = statusCodeHandler(resp)
	if err != nil {
		return GetAdsResponse{}, fmt.Errorf("statusCodeHandler: %w", err)
	}

	var errResp APIError

	err = json.Unmarshal(body, &errResp)
	if err != nil {
		return GetAdsResponse{}, err
	}

	if errResp.Err.ErrorCode != 0 {
		return GetAdsResponse{}, errResp
	}

	var response GetAdsResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return GetAdsResponse{}, err
	}

	return response, nil
}

func (c *Client) createGetAdsRequest(ctx context.Context, params GetAdsRequest) (*http.Request, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.direct.yandex.com/json/v5/ads", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	c.buildHeader(req)

	return req, nil
}

type GetAdsRequest struct {
	Method string           `json:"method"`
	Params GetAdsDefinition `json:"params"`
}

type GetAdsDefinition struct {
	Selection        AdsSelectionCriteria `json:"SelectionCriteria,omitempty"`
	FieldNames       []string             `json:"FieldNames"`
	TextAdFieldNames []string             `json:"TextAdFieldNames,omitempty"`
}

type AdsSelectionCriteria struct {
	IDs                         []int    `json:"Ids,omitempty"`
	States                      []string `json:"States,omitempty"`
	Statuses                    []string `json:"Statuses,omitempty"`
	CampaignIDs                 []int    `json:"CampaignIds,omitempty"`
	AdGroupIDs                  []int    `json:"AdGroupIds,omitempty"`
	Types                       []string `json:"Types,omitempty"`
	Mobile                      *string  `json:"Mobile,omitempty"`
	VCardIDs                    []int    `json:"VCardIds,omitempty"`
	SiteLinkSetIDs              []int    `json:"SitelinkSetIds,omitempty"`
	AdImageHashes               []string `json:"AdImageHashes,omitempty"`
	VCardModerationStatuses     []string `json:"VCardModerationStatuses,omitempty"`
	SiteLinksModerationStatuses []string `json:"SitelinksModerationStatuses,omitempty"`
	AdImageModerationStatuses   []string `json:"AdImageModerationStatuses,omitempty"`
	AdExtensionIDs              []int    `json:"AdExtensionIds,omitempty"`
}

type GetAdsResponse struct {
	Result struct {
		Ads []AdItem `json:"Ads"`
	} `json:"result"`
}

type AdItem struct {
	ID         int    `json:"Id"`
	CampaignID int    `json:"CampaignId"`
	AdGroupID  int64  `json:"AdGroupId"`
	TextAd     string `json:"TextAd"`
}

func (r *GetAdsResponse) UnmarshalJSON(data []byte) error {
	var in struct {
		Result struct {
			Ads []struct {
				ID         int             `json:"Id"`
				CampaignID int             `json:"CampaignId"`
				AdGroupID  int64           `json:"AdGroupId"`
				TextAd     json.RawMessage `json:"TextAd"`
			} `json:"Ads"`
		} `json:"result"`
	}
	if err := json.Unmarshal(data, &in); err != nil {
		return err
	}

	r.Result.Ads = make([]AdItem, 0, len(in.Result.Ads))
	for _, ad := range in.Result.Ads {
		r.Result.Ads = append(r.Result.Ads, AdItem{
			ID:         ad.ID,
			CampaignID: ad.CampaignID,
			AdGroupID:  ad.AdGroupID,
			TextAd:     string(ad.TextAd),
		})
	}

	return nil
}
