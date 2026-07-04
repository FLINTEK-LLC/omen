package ingest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/FLINTEK-LLC/omen/internal/store"
)

// RansomwareLiveClient is a Source backed by the public Ransomware.live v2
// API (https://www.ransomware.live/apidocs). It requires no API key.
type RansomwareLiveClient struct {
	baseURL string
	http    *http.Client
}

// NewRansomwareLiveClient builds a client against baseURL (e.g.
// "https://api.ransomware.live/v2").
func NewRansomwareLiveClient(baseURL string) *RansomwareLiveClient {
	return &RansomwareLiveClient{
		baseURL: baseURL,
		http:    &http.Client{Timeout: 30 * time.Second},
	}
}

// FetchRecentVictims retrieves the latest disclosed victims from
// /recentvictims.
func (c *RansomwareLiveClient) FetchRecentVictims(ctx context.Context) ([]store.Victim, error) {
	var raw []map[string]any
	if err := c.getJSON(ctx, "/recentvictims", &raw); err != nil {
		return nil, fmt.Errorf("fetch recent victims: %w", err)
	}

	victims := make([]store.Victim, 0, len(raw))
	for _, r := range raw {
		victims = append(victims, victimFromRaw(r))
	}
	return victims, nil
}

// FetchVictimsByMonth retrieves all victims disclosed in the given
// year/month from /victims/{year}/{month}. Unlike FetchRecentVictims (a
// rolling window capped at 100 records), this endpoint returns the full set
// for that month, so it's the source used for historical backfill.
func (c *RansomwareLiveClient) FetchVictimsByMonth(ctx context.Context, year int, month time.Month) ([]store.Victim, error) {
	var raw []map[string]any
	path := fmt.Sprintf("/victims/%d/%d", year, int(month))
	if err := c.getJSON(ctx, path, &raw); err != nil {
		return nil, fmt.Errorf("fetch victims for %d-%02d: %w", year, month, err)
	}

	victims := make([]store.Victim, 0, len(raw))
	for _, r := range raw {
		victims = append(victims, victimFromRaw(r))
	}
	return victims, nil
}

// FetchGroups retrieves the full list of tracked ransomware groups from
// /groups.
func (c *RansomwareLiveClient) FetchGroups(ctx context.Context) ([]store.Group, error) {
	var raw []map[string]any
	if err := c.getJSON(ctx, "/groups", &raw); err != nil {
		return nil, fmt.Errorf("fetch groups: %w", err)
	}

	groups := make([]store.Group, 0, len(raw))
	for _, r := range raw {
		name, _ := r["name"].(string)
		if name == "" {
			continue
		}
		groups = append(groups, store.Group{Name: name, ProfileJSON: r})
	}
	return groups, nil
}

func (c *RansomwareLiveClient) getJSON(ctx context.Context, path string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d from %s", resp.StatusCode, path)
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

// victimFromRaw maps a raw Ransomware.live victim record (see field names in
// the API docs: victim, group, attackdate, country, domain, infostealer,
// press, updates) onto our store model, keeping the full original payload in
// RawSource for audit.
//
// Ransomware.live supplies domain directly for many victims; when it
// doesn't, Domain is left blank. A name-to-domain guesser or manual
// override for those remaining cases is deferred to a follow-up (see spec
// section 12).
func victimFromRaw(r map[string]any) store.Victim {
	name, _ := r["victim"].(string)
	group, _ := r["group"].(string)
	date, _ := r["attackdate"].(string)
	country, _ := r["country"].(string)
	domain, _ := r["domain"].(string)

	var press []string
	if arr, ok := r["press"].([]any); ok {
		for _, item := range arr {
			switch v := item.(type) {
			case string:
				if v != "" {
					press = append(press, v)
				}
			case map[string]any:
				if link, ok := v["link"].(string); ok && link != "" {
					press = append(press, link)
				} else if url, ok := v["url"].(string); ok && url != "" {
					press = append(press, url)
				}
			}
		}
	}

	infostealer, _ := r["infostealer"].(map[string]any)
	if infostealer == nil {
		infostealer = map[string]any{}
	}

	return store.Victim{
		ID:          store.VictimID(name, group, date),
		VictimName:  name,
		GroupName:   group,
		AttackDate:  date,
		Country:     country,
		Domain:      domain,
		PressLinks:  press,
		Infostealer: infostealer,
		RawSource:   r,
	}
}
