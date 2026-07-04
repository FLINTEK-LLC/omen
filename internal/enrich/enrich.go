// Package enrich adds Shodan exposure data and CISA KEV known-exploited
// status to newly ingested victims.
package enrich

import "context"

// ShodanResult is the subset of a Shodan host lookup OMEN persists.
type ShodanResult struct {
	IP    string
	Ports []int
	CVEs  []string
}

// Enricher looks up exposure data for a victim's domain or IP. NoopEnricher
// is used when no Shodan API key is configured.
type Enricher interface {
	Lookup(ctx context.Context, target string) (*ShodanResult, error)
}

// NoopEnricher is a no-op Enricher used when Shodan enrichment is disabled
// (no API key configured). It lets the rest of the ingest/enrich pipeline
// run unchanged.
type NoopEnricher struct{}

// Lookup always returns (nil, nil), meaning "no result, not an error".
func (NoopEnricher) Lookup(ctx context.Context, target string) (*ShodanResult, error) {
	return nil, nil
}
