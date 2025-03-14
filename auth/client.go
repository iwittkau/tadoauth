package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var ErrTodoAPI = errors.New("tado API error")

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) RequestRegistration(ctx context.Context) (*DeviceAuthorizeResponse, error) {
	form := make(url.Values, 2)
	form.Add("client_id", tadoClientID)
	form.Add("scope", clientScope)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, deviceAuthorizeURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("creating device authorization request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("device authorization request: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("%w: %s - %s", ErrTodoAPI, res.Status, string(data))
	}

	var deviceAuthRes DeviceAuthorizeResponse
	if err := json.NewDecoder(res.Body).Decode(&deviceAuthRes); err != nil {
		return nil, fmt.Errorf("unmarshaling device authorize response: %w", err)
	}
	return &deviceAuthRes, nil
}

func (c *Client) ExchangeDeviceCode(ctx context.Context, deviceCode string) (*TokenResponse, error) {
	form := make(url.Values, 2)
	form.Add("client_id", tadoClientID)
	form.Add("grant_type", grantTypeDeviceCode)
	form.Add("device_code", deviceCode)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("creating token request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token request: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("%w: %s - %s", ErrTodoAPI, res.Status, string(data))
	}

	var tokenRes TokenResponse
	if err := json.NewDecoder(res.Body).Decode(&tokenRes); err != nil {
		return nil, fmt.Errorf("unmarshaling token response: %w", err)
	}

	return &tokenRes, nil
}

func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	form := make(url.Values, 2)
	form.Add("client_id", tadoClientID)
	form.Add("grant_type", "refresh_token")
	form.Add("refresh_token", refreshToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("creating token request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token request: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("%w: %s - %s", ErrTodoAPI, res.Status, string(data))
	}

	var tokenRes TokenResponse
	if err := json.NewDecoder(res.Body).Decode(&tokenRes); err != nil {
		return nil, fmt.Errorf("unmarshaling token response: %w", err)
	}

	return &tokenRes, nil
}
