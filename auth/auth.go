package auth

import (
	"log"
	"time"
)

const (
	deviceAuthorizeURL    = "https://login.tado.com/oauth2/device_authorize"
	tokenURL              = "https://login.tado.com/oauth2/token"
	tadoClientID          = "1bb50063-6b0c-4d11-bd99-387f4a91cc46"
	clientScope           = "offline_access"
	grantTypeDeviceCode   = "urn:ietf:params:oauth:grant-type:device_code"
	grantTypeRefreshTOken = "refresh_token"
)

type DeviceAuthorizeResponse struct {
	DeviceCode              string `json:"device_code"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
}

func (r DeviceAuthorizeResponse) ExpiresAt() time.Time {
	return time.Now().Add(time.Duration(r.ExpiresIn) * time.Second)
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	UserID       string `json:"userId"`
}

func (r TokenResponse) Print() {
	log.Println("AccessToken: ", r.AccessToken)
	log.Println("ExpiresIn:   ", r.ExpiresIn)
	log.Println("RefreshToken:", r.RefreshToken)
	log.Println("Scope:       ", r.Scope)
	log.Println("TokenType:   ", r.TokenType)
	log.Println("UserID:      ", r.UserID)
}
