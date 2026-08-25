package email

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"
)

var ErrNotConfigured = errors.New("email service is not configured")

type Service interface {
	SendVerificationCode(context.Context, string, string, string) error
}

type resendService struct {
	apiKey string
	from   string
	client *http.Client
}

type resendEmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

func NewResendService(apiKey, from string) Service {
	return &resendService{apiKey: strings.TrimSpace(apiKey), from: strings.TrimSpace(from), client: &http.Client{Timeout: 10 * time.Second}}
}

func (s *resendService) SendVerificationCode(ctx context.Context, recipient, name, code string) error {
	if s.apiKey == "" || s.from == "" {
		return ErrNotConfigured
	}

	recipient = strings.TrimSpace(recipient)
	name = html.EscapeString(strings.TrimSpace(name))
	code = html.EscapeString(code)
	payload, err := json.Marshal(resendEmailRequest{
		From: s.from, To: []string{recipient}, Subject: "Verify your GoTask email",
		HTML: fmt.Sprintf("<p>Hello %s,</p><p>Your GoTask verification code is:</p><p><strong>%s</strong></p><p>This code expires in 10 minutes.</p>", name, code),
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.resend.com/emails", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "gotask/1.0")
	res, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("resend returned status %d", res.StatusCode)
	}
	return nil
}

func GenerateVerificationCode() (string, error) {
	var digits [6]byte
	if _, err := rand.Read(digits[:]); err != nil {
		return "", err
	}
	for i := range digits {
		digits[i] = '0' + digits[i]%10
	}
	return string(digits[:]), nil
}
