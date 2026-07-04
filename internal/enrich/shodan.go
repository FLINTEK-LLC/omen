package enrich

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

// ShodanClient is an Enricher backed by the real Shodan host lookup API
// (https://developer.shodan.io/api). Calls are serialized through a ticker
// so callers don't need to reason about rate limits themselves; on a 429 the
// caller's retry (if any) will simply wait for the next tick.
type ShodanClient struct {
	apiKey   string
	http     *http.Client
	throttle <-chan time.Time
}

// NewShodanClient builds a client that issues at most one request per
// minInterval.
func NewShodanClient(apiKey string, minInterval time.Duration) *ShodanClient {
	return &ShodanClient{
		apiKey:   apiKey,
		http:     &http.Client{Timeout: 15 * time.Second},
		throttle: time.Tick(minInterval),
	}
}

type shodanHostResponse struct {
	IPStr string `json:"ip_str"`
	Data  []struct {
		Port  int             `json:"port"`
		Vulns json.RawMessage `json:"vulns"` // Shodan returns an object keyed by CVE ID
	} `json:"data"`
}

// Lookup resolves target (an IP or hostname) to an IP if needed, then
// queries Shodan's host endpoint for open ports and any flagged CVEs.
func (c *ShodanClient) Lookup(ctx context.Context, target string) (*ShodanResult, error) {
	ip := target
	if net.ParseIP(target) == nil {
		addrs, err := net.DefaultResolver.LookupHost(ctx, target)
		if err != nil {
			return nil, fmt.Errorf("resolve %s: %w", target, err)
		}
		if len(addrs) == 0 {
			return nil, fmt.Errorf("no addresses found for %s", target)
		}
		ip = addrs[0]
	}

	select {
	case <-c.throttle:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	url := fmt.Sprintf("https://api.shodan.io/shodan/host/%s?key=%s", ip, c.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// No Shodan data for this host; not an error.
		return &ShodanResult{IP: ip}, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("shodan returned status %d", resp.StatusCode)
	}

	var parsed shodanHostResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("decode shodan response: %w", err)
	}

	result := &ShodanResult{IP: ip}
	cveSet := map[string]struct{}{}
	for _, d := range parsed.Data {
		result.Ports = append(result.Ports, d.Port)
		if len(d.Vulns) == 0 {
			continue
		}
		var vulnsByCVE map[string]json.RawMessage
		if err := json.Unmarshal(d.Vulns, &vulnsByCVE); err == nil {
			for cve := range vulnsByCVE {
				cveSet[cve] = struct{}{}
			}
		}
	}
	for cve := range cveSet {
		result.CVEs = append(result.CVEs, cve)
	}

	return result, nil
}
