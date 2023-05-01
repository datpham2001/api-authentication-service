package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"realworld-authentication/config"
	"time"
)

type GoogleOauthToken struct {
	AccessToken string
	TokenID     string
}

type GoogleUserInfo struct {
	ID    string
	Email string
	Name  string
}

func GetGoogleOauthToken(code string) (*GoogleOauthToken, error) {
	const rootUrl = "https://oauth2.googleapis.com/token"

	values := url.Values{}
	values.Add("grant_type", "authorization_code")
	values.Add("code", code)
	values.Add("client_id", config.AppConfig.GoogleOauthClientID)
	values.Add("client_secret", config.AppConfig.GoogleOauthSecret)
	values.Add("redirect_uri", config.AppConfig.GoogleOauthRedirectUrl)

	query := values.Encode()

	req, err := http.NewRequest("POST", rootUrl, bytes.NewBufferString(query))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve the token")
	}

	var resBody bytes.Buffer
	_, err = io.Copy(&resBody, res.Body)
	if err != nil {
		return nil, err
	}

	var googleOauthTokenResp map[string]interface{}
	err = json.Unmarshal(resBody.Bytes(), &googleOauthTokenResp)
	if err != nil {
		return nil, err
	}

	tokenBody := &GoogleOauthToken{
		AccessToken: googleOauthTokenResp["access_token"].(string),
		TokenID:     googleOauthTokenResp["id_token"].(string),
	}

	return tokenBody, nil
}

func GetGoogleUserInfo(accessToken, tokenID string) (*GoogleUserInfo, error) {
	rootUrl := fmt.Sprintf("https://www.googleapis.com/oauth2/v1/userinfo?alt=json&access_token=%s", accessToken)

	req, err := http.NewRequest("GET", rootUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer &s", tokenID))
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve user")
	}

	var resBody bytes.Buffer
	_, err = io.Copy(&resBody, res.Body)
	if err != nil {
		return nil, err
	}

	var googleUserInfoResp map[string]interface{}
	err = json.Unmarshal(resBody.Bytes(), &googleUserInfoResp)
	if err != nil {
		return nil, err
	}

	userInfo := &GoogleUserInfo{
		ID:    googleUserInfoResp["id"].(string),
		Email: googleUserInfoResp["email"].(string),
		Name:  googleUserInfoResp["name"].(string),
	}

	return userInfo, nil
}
