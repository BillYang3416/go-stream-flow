package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/bgg/go-flow-gateway/internal/usecase/dto"
	"github.com/bgg/go-flow-gateway/pkg/logger"
)

type LineTokenService struct {
	logger logger.Logger
}

func NewLineTokenService(l logger.Logger) *LineTokenService {
	return &LineTokenService{logger: l}
}

func (l *LineTokenService) VerifyIDToken(idToken, clientID string) (*dto.LineUserProfile, error) {
	data := fmt.Sprintf("id_token=%s&client_id=%s", idToken, clientID)

	req, err := http.NewRequest("POST", "https://api.line.me/oauth2/v2.1/verify", bytes.NewBufferString(data))
	if err != nil {
		return nil, fmt.Errorf("error sending request to LINE: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error getting response from LINE: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from LINE: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body from LINE: %v", err)
	}

	var profile dto.LineUserProfile
	if err := json.Unmarshal(body, &profile); err != nil {
		return nil, fmt.Errorf("error unmarshalling response body from LINE: %v", err)
	}

	return &profile, nil
}

func (l *LineTokenService) ExchangeCodeForTokens(code, domainUrl string) (*dto.TokenResponse, error) {

	tokenEndpoint := "https://api.line.me/oauth2/v2.1/token"

	data := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {fmt.Sprintf("%s/api/v1/auth/line-callback", domainUrl)},
		"client_id":     {os.Getenv("LINE_CHANNEL_ID")},
		"client_secret": {os.Getenv("LINE_CHANNEL_SECRET")},
	}

	req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tr dto.TokenResponse
	if err := json.Unmarshal(body, &tr); err != nil {
		return nil, err
	}

	return &tr, nil
}
