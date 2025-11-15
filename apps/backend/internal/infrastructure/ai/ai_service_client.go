// apps/backend/internal/infrastructure/ai/ai_service_client.go

package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AnalyzeRequest ساختار درخواستی است که به سرویس AI ارسال می‌شود.
type AnalyzeRequest struct {
	Text    string `json:"text"`
	Context string `json:"context,omitempty"`
}

// AnalyzeResponse ساختار پاسخی است که از سرویس AI دریافت می‌شود.
type AnalyzeResponse struct {
	Analysis string `json:"analysis"`
}

// Client کلاینت برای ارتباط با سرویس AI است.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient یک نمونه جدید از AI Service Client ایجاد می‌کند.
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 60 * time.Second, // Timeout طولانی‌تر برای پردازش‌های AI
		},
	}
}

// AnalyzeText یک متن را برای تحلیل به سرویس AI ارسال می‌کند.
func (c *Client) AnalyzeText(ctx context.Context, text, context string) (*AnalyzeResponse, error) {
	// 1. آماده‌سازی درخواست JSON
	reqPayload := AnalyzeRequest{
		Text:    text,
		Context: context,
	}
	payloadBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %w", err)
	}

	// 2. ساخت درخواست HTTP
	url := fmt.Sprintf("%s/analyze", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	// TODO: Add internal service-to-service authentication (e.g., a shared secret header)

	// 3. ارسال درخواست
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call ai service: %w", err)
	}
	defer resp.Body.Close()

	// 4. بررسی پاسخ
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ai service returned non-200 status: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	// 5. خواندن و بازگرداندن پاسخ
	var analyzeResp AnalyzeResponse
	if err := json.NewDecoder(resp.Body).Decode(&analyzeResp); err != nil {
		return nil, fmt.Errorf("failed to decode ai service response: %w", err)
	}

	if analyzeResp.Analysis == "" {
		return nil, errors.New("ai service returned empty analysis")
	}

	return &analyzeResp, nil
}