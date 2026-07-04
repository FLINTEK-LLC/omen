package watchlist

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookNotifier POSTs a JSON payload to a webhook URL (Slack, Teams, n8n,
// or any endpoint that accepts a JSON body), per spec section 8's
// recommendation that webhook be the lowest-effort first integration point.
type WebhookNotifier struct {
	http       *http.Client
	defaultURL string
}

// NewWebhookNotifier builds a notifier that falls back to defaultURL when a
// watchlist entry doesn't specify its own notify_via.
func NewWebhookNotifier(defaultURL string) *WebhookNotifier {
	return &WebhookNotifier{
		http:       &http.Client{Timeout: 10 * time.Second},
		defaultURL: defaultURL,
	}
}

// Notify POSTs event as JSON to notifyVia, or defaultURL if notifyVia is
// empty. It is a no-op if neither is set.
func (w *WebhookNotifier) Notify(ctx context.Context, notifyVia string, event Event) error {
	url := notifyVia
	if url == "" {
		url = w.defaultURL
	}
	if url == "" {
		return nil
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.http.Do(req)
	if err != nil {
		return fmt.Errorf("post webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}
	return nil
}
