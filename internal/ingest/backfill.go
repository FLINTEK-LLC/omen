package ingest

import (
	"context"
	"log"
	"time"

	"github.com/FLINTEK-LLC/omen/internal/store"
)

// rateLimitCooldown matches Ransomware.live's documented "1 request per
// minute per endpoint" limit for its free API tier, with a small margin.
const rateLimitCooldown = 65 * time.Second

// BackfillMonths pulls full monthly victim data from
// /victims/{year}/{month} for the current month and the (months-1)
// preceding it, upserting everything into the store. It's a one-time
// startup job, not the ongoing poll loop: it only persists victims and does
// not trigger enrichment or watchlist notifications, since re-alerting on
// months-old disclosures every time the process restarts would be noise.
// Re-running it is safe -- already-seen victims are simply re-upserted.
func BackfillMonths(ctx context.Context, client *RansomwareLiveClient, st *store.Store, months int) {
	if months <= 0 {
		return
	}

	now := time.Now().UTC()
	for i := 0; i < months; i++ {
		if ctx.Err() != nil {
			return
		}

		t := now.AddDate(0, -i, 0)
		victims, err := client.FetchVictimsByMonth(ctx, t.Year(), t.Month())
		if err != nil {
			log.Printf("ingest: backfill %d-%02d: %v", t.Year(), int(t.Month()), err)
		} else {
			upserted := 0
			for _, v := range victims {
				if err := st.UpsertVictim(ctx, v); err != nil {
					log.Printf("ingest: backfill upsert victim %s: %v", v.ID, err)
					continue
				}
				upserted++
			}
			log.Printf("ingest: backfill %d-%02d: upserted %d of %d victims", t.Year(), int(t.Month()), upserted, len(victims))
		}

		if i < months-1 {
			select {
			case <-ctx.Done():
				return
			case <-time.After(rateLimitCooldown):
			}
		}
	}

	log.Println("ingest: backfill complete")
}
