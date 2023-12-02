package token

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// UserProfile represents the user's profile information returned by LINE
type LineUserProfile struct {
	Iss     string `json:"iss"`
	Sub     string `json:"sub"`
	Aud     string `json:"aud"`
	Exp     int64  `json:"exp"`
	Iat     int64  `json:"iat"`
	Nonce   string `json:"nonce"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	Email   string `json:"email"`
}

func VerifyLineIDToken(idToken, clientId string) (*LineUserProfile, error) {
	data := fmt.Sprintf("id_token=%s&client_id=%s", idToken, clientId)

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

	var profile LineUserProfile
	if err := json.Unmarshal(body, &profile); err != nil {
		return nil, fmt.Errorf("error unmarshalling response body from LINE: %v", err)
	}

	return &profile, nil
}
