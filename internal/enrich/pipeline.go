package enrich

import (
	"context"
	"fmt"
	"time"

	"github.com/FLINTEK-LLC/omen/internal/store"
)

// EnrichVictim looks up exposure data for v (if it has a resolvable domain),
// persists the Shodan snapshot, and cross-references any flagged CVEs
// against the KEV catalog, persisting matches. It is a no-op if v has no
// domain or the enricher returns no result.
func EnrichVictim(ctx context.Context, st *store.Store, enricher Enricher, kev *KEVCatalog, v store.Victim) error {
	if v.Domain == "" {
		return nil
	}

	result, err := enricher.Lookup(ctx, v.Domain)
	if err != nil {
		return fmt.Errorf("shodan lookup for %s: %w", v.Domain, err)
	}
	if result == nil {
		return nil
	}

	now := time.Now().UTC()
	if err := st.InsertShodanSnapshot(ctx, store.ShodanSnapshot{
		VictimID:  v.ID,
		IP:        result.IP,
		Ports:     result.Ports,
		CVEs:      result.CVEs,
		QueriedAt: now,
	}); err != nil {
		return fmt.Errorf("store shodan snapshot: %w", err)
	}

	for _, cve := range result.CVEs {
		entry, ok := kev.Lookup(cve)
		if !ok {
			continue
		}
		if err := st.InsertKEVMatch(ctx, store.KEVMatch{
			VictimID:  v.ID,
			CVEID:     entry.CVEID,
			KEVAdded:  entry.DateAdded,
			MatchedAt: now,
		}); err != nil {
			return fmt.Errorf("store kev match for %s: %w", cve, err)
		}
	}

	return nil
}
