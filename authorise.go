package direct_sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type AuthByCode struct {
	DeviceCode      string `json:"device_code"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
	UserCode        string `json:"user_code"`
	VerificationURL string `json:"verification_url"`
}

type responseError struct {
	Err string `json:"error"`
	Msg string `json:"error_description"`
}

func (e responseError) Error() string {
	return fmt.Sprintf("%s: %s", e.Err, e.Msg)
}

var ErrAuthorisationPending = responseError{
	Err: "authorization_pending",
	Msg: "User has not yet authorized your application",
}

type Token struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

func (c *Client) Authorise(message chan<- string) error {
	authData, err := c.GetAccessCode()
	if err != nil {
		return err
	}

	message <- fmt.Sprintf("user code: %s | %s \n", authData.UserCode, authData.VerificationURL)
	token := Token{}
	done := make(chan bool, 1)
	duration := time.Duration(authData.Interval)*time.Second + 20*time.Second
	ticker := time.NewTicker(duration)

	defer ticker.Stop()

outerLoop:
	for {
		select {
		case <-done:
			message <- "nil"
			close(message)
			return nil
		case <-ticker.C:
			token, err = c.GetTokenByCode(authData)
			if errors.Is(err, ErrAuthorisationPending) {
				fmt.Println(err)

				continue
			} else if err != nil {
				return err
			}
			if token.AccessToken != "" {
				done <- true

				break outerLoop
			}
		}
	}
	message <- "token: " + token.AccessToken
	close(message)

	return nil
}

func (c *Client) GetAccessCode() (auth AuthByCode, err error) {
	host := "oauth.yandex.ru"

	reqAccessURL := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "/device/code",
	}
	param := url.Values{}
	param.Add("client_id", c.App.ID)

	resp, err := c.Tr.PostForm(reqAccessURL.String(), param)
	if err != nil {
		return auth, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return auth, fmt.Errorf("StatusCode: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return auth, err
	}

	err = json.Unmarshal(body, &auth)
	if err != nil {
		return auth, err
	}

	return auth, err
}

func (c *Client) GetTokenByCode(code AuthByCode) (token Token, err error) {
	host := "oauth.yandex.ru"

	reqAccessURL := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "token",
	}
	param := url.Values{
		"client_id":     {c.App.ID},
		"grant_type":    {"device_code"},
		"code":          {code.DeviceCode},
		"client_secret": {c.App.Secret},
	}

	resp, err := c.Tr.PostForm(reqAccessURL.String(), param)
	if err != nil {
		return token, err
	}
	defer resp.Body.Close()

	//if resp.StatusCode != http.StatusOK {
	//	return token, fmt.Errorf("StatusCode: %v", resp.StatusCode)
	//}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return token, err
	}

	err = json.Unmarshal(body, &token)
	if err != nil {
		return token, err
	}

	respErr := responseError{}

	err = json.Unmarshal(body, &respErr)
	if err != nil {
		return token, err
	}

	if respErr.Err != "" {
		return token, respErr
	}

	return token, err
}
