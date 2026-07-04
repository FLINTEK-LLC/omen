package enrich

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// KEVEntry is one entry from the CISA Known Exploited Vulnerabilities
// catalog.
type KEVEntry struct {
	CVEID             string `json:"cveID"`
	VulnerabilityName string `json:"vulnerabilityName"`
	DateAdded         string `json:"dateAdded"`
	ShortDescription  string `json:"shortDescription"`
}

type kevFeed struct {
	Vulnerabilities []KEVEntry `json:"vulnerabilities"`
}

// KEVCatalog holds an in-memory, disk-cached copy of the CISA KEV catalog,
// refreshed periodically from feedURL.
type KEVCatalog struct {
	feedURL   string
	cachePath string
	http      *http.Client

	mu    sync.RWMutex
	byCVE map[string]KEVEntry
}

// NewKEVCatalog builds a catalog backed by feedURL and cached to disk at
// cachePath.
func NewKEVCatalog(feedURL, cachePath string) *KEVCatalog {
	return &KEVCatalog{
		feedURL:   feedURL,
		cachePath: cachePath,
		http:      &http.Client{Timeout: 30 * time.Second},
		byCVE:     map[string]KEVEntry{},
	}
}

// Lookup returns the KEV entry for cveID, if the CVE is known-exploited.
func (k *KEVCatalog) Lookup(cveID string) (KEVEntry, bool) {
	k.mu.RLock()
	defer k.mu.RUnlock()
	e, ok := k.byCVE[strings.ToUpper(cveID)]
	return e, ok
}

// LoadFromDisk populates the catalog from the on-disk cache, if present. Call
// this at startup before the first network refresh completes, so lookups
// aren't empty while waiting on the network.
func (k *KEVCatalog) LoadFromDisk() error {
	data, err := os.ReadFile(k.cachePath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read kev cache: %w", err)
	}
	return k.load(data)
}

// Refresh fetches the latest catalog from feedURL, updates the in-memory
// index, and writes the result to the on-disk cache.
func (k *KEVCatalog) Refresh(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, k.feedURL, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	resp, err := k.http.Do(req)
	if err != nil {
		return fmt.Errorf("fetch kev feed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("kev feed returned status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read kev feed body: %w", err)
	}

	if err := k.load(data); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(k.cachePath), 0o755); err != nil {
		return fmt.Errorf("create cache dir: %w", err)
	}
	if err := os.WriteFile(k.cachePath, data, 0o644); err != nil {
		return fmt.Errorf("write kev cache: %w", err)
	}
	return nil
}

func (k *KEVCatalog) load(data []byte) error {
	var feed kevFeed
	if err := json.Unmarshal(data, &feed); err != nil {
		return fmt.Errorf("parse kev feed: %w", err)
	}

	byCVE := make(map[string]KEVEntry, len(feed.Vulnerabilities))
	for _, e := range feed.Vulnerabilities {
		byCVE[strings.ToUpper(e.CVEID)] = e
	}

	k.mu.Lock()
	k.byCVE = byCVE
	k.mu.Unlock()
	return nil
}
