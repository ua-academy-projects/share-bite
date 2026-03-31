package email

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

const (
	resendSendEmailURL = "https://api.resend.com/emails"
	resetEmailSubject  = "Reset your Share Bite password"
)

type Sender interface {
	SendPasswordResetToken(ctx context.Context, toEmail, token string) error
}

type resendSender struct {
	apiKey    string
	fromEmail string
	client    *http.Client
}

type resendSendEmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

func NewResendSender(apiKey, fromEmail string) (Sender, error) {
	if apiKey == "" {
		return nil, errors.New("resend api key is empty")
	}
	if fromEmail == "" {
		return nil, errors.New("resend from email is empty")
	}

	return &resendSender{
		apiKey:    apiKey,
		fromEmail: fromEmail,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

func (s *resendSender) SendPasswordResetToken(ctx context.Context, toEmail, token string) error {
	if toEmail == "" {
		return errors.New("recipient email is empty")
	}
	if token == "" {
		return errors.New("reset token is empty")
	}

	logger.InfoKV(ctx, "sending password reset email")

	reqBody, err := json.Marshal(resendSendEmailRequest{
		From:    s.fromEmail,
		To:      []string{toEmail},
		Subject: resetEmailSubject,
		HTML: fmt.Sprintf(
			"<p>Hello,</p><p>You requested a password reset for your Share Bite account.</p><p>Your reset token:</p><p><strong>%s</strong></p><p>Use this token with the reset password API endpoint.</p>",
			token,
		),
	})
	if err != nil {
		return fmt.Errorf("marshal resend request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, resendSendEmailURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("build resend request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("send resend request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("resend send email failed: status=%d body=%s", resp.StatusCode, string(respBody))
	}

	logger.InfoKV(ctx, "password reset email sent")
	return nil
}
