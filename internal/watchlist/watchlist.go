// Package watchlist matches newly ingested victims against user-configured
// patterns and dispatches notifications on match.
package watchlist

import (
	"context"
	"strings"

	"github.com/FLINTEK-LLC/omen/internal/store"
)

// Event describes a watchlist match, passed to a Notifier.
type Event struct {
	VictimID       string `json:"victim_id"`
	VictimName     string `json:"victim_name"`
	GroupName      string `json:"group_name"`
	Domain         string `json:"domain,omitempty"`
	MatchedPattern string `json:"matched_pattern"`
	MatchedLabel   string `json:"matched_label,omitempty"`
}

// Notifier dispatches a watchlist match event to a notification channel
// (webhook, email, etc). notifyVia is the entry's configured destination,
// which may be empty to fall back to a Notifier-specific default.
type Notifier interface {
	Notify(ctx context.Context, notifyVia string, event Event) error
}

// Match returns every watchlist entry whose pattern is a case-insensitive
// substring of the victim's name or domain. Kept intentionally dumb per spec
// section 8 -- revisit with fuzzy matching only if false positive/negative
// rate becomes a problem in practice.
func Match(entries []store.WatchlistEntry, victimName, domain string) []store.WatchlistEntry {
	nameLower := strings.ToLower(victimName)
	domainLower := strings.ToLower(domain)

	var matches []store.WatchlistEntry
	for _, e := range entries {
		pattern := strings.ToLower(strings.TrimSpace(e.Pattern))
		if pattern == "" {
			continue
		}
		if strings.Contains(nameLower, pattern) || (domainLower != "" && strings.Contains(domainLower, pattern)) {
			matches = append(matches, e)
		}
	}
	return matches
}
