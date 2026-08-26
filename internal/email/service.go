package email

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gotask/internal/email/template"
	emailtypes "gotask/internal/types/email"
)

var ErrNotConfigured = errors.New("email service is not configured")

type Service interface {
	SendVerificationCode(context.Context, string, string, string) error
	SendEmailChangedNotification(context.Context, string, string, string) error
}

type resendService struct {
	apiKey string
	from   string
	client *http.Client
}

func NewResendService(apiKey, from string) Service {
	return &resendService{apiKey: strings.TrimSpace(apiKey), from: strings.TrimSpace(from), client: &http.Client{Timeout: 10 * time.Second}}
}

func (s *resendService) SendVerificationCode(ctx context.Context, recipient, name, code string) error {
	return s.send(ctx, recipient, "Verify your GoTask email", template.VerificationCodeTemplate(name, code))
}

func (s *resendService) SendEmailChangedNotification(ctx context.Context, recipient, name, newEmail string) error {
	return s.send(ctx, recipient, "Your GoTask email was changed", template.EmailChangedTemplate(name, newEmail))
}

func (s *resendService) send(ctx context.Context, recipient, subject, body string) error {
	if s.apiKey == "" || s.from == "" {
		return ErrNotConfigured
	}

	recipient = strings.TrimSpace(recipient)
	payload, err := json.Marshal(emailtypes.ResendEmailRequest{
		From: s.from, To: []string{recipient}, Subject: subject, HTML: body,
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
