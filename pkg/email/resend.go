package email

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

const (
	resendSendEmailURL = "https://api.resend.com/emails"
	resetEmailSubject  = "Reset your Share Bite password"
	templateResetPath  = "templates/password_reset.html"
)

//go:embed templates/*
var templateFS embed.FS

type resendSender struct {
	apiKey    string
	fromEmail string
	client    *http.Client
	tpl       *template.Template
}

type passwordResetTemplateData struct {
	Token string
}

type resendSendEmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

func NewResendSender(apiKey, fromEmail string) Sender {
	return &resendSender{
		apiKey:    apiKey,
		fromEmail: fromEmail,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *resendSender) SendPasswordResetToken(ctx context.Context, toEmail, token string) error {
	if toEmail == "" {
		return errors.New("recipient email is empty")
	}
	if token == "" {
		return errors.New("reset token is empty")
	}

	return s.SendEmail(ctx, toEmail, resetEmailSubject, "password_reset", map[string]any{
		"token": token,
	})
}

func (s *resendSender) SendEmail(ctx context.Context, toEmail, subject, templateName string, data map[string]any) error {
	if toEmail == "" {
		return errors.New("recipient email is empty")
	}
	if subject == "" {
		return errors.New("subject is empty")
	}
	if templateName == "" {
		return errors.New("template name is empty")
	}

	logger.InfoKV(ctx, "sending email via Resend", "to", toEmail, "subject", subject, "template", templateName)

	templatePath := fmt.Sprintf("templates/%s.html", templateName)
	tpl, err := template.ParseFS(templateFS, templatePath)
	if err != nil {
		return fmt.Errorf("parse email template %s: %w", templatePath, err)
	}

	var htmlBody bytes.Buffer
	if err := tpl.Execute(&htmlBody, data); err != nil {
		return fmt.Errorf("render email template %s: %w", templatePath, err)
	}

	reqBody, err := json.Marshal(resendSendEmailRequest{
		From:    s.fromEmail,
		To:      []string{toEmail},
		Subject: subject,
		HTML:    htmlBody.String(),
	})
	if err != nil {
		return fmt.Errorf("marshal resend request: %w", err)
	}

	// Send request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, resendSendEmailURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("build resend request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		logger.ErrorKV(ctx, "email send failed", "to", toEmail, "error", err.Error())
		return fmt.Errorf("send resend request: %w", err)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			logger.ErrorKV(ctx, "email send failed", "to", toEmail, "status", resp.StatusCode, "error", readErr.Error())
			return fmt.Errorf("resend send email failed: status=%d, read body: %w", resp.StatusCode, readErr)
		}

		logger.ErrorKV(ctx, "email send failed", "to", toEmail, "status", resp.StatusCode, "body", string(body))
		return fmt.Errorf("resend send email failed: status=%d, body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	logger.InfoKV(ctx, "email sent successfully", "to", toEmail, "subject", subject)
	return nil
}
