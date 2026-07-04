// Package ingest polls Ransomware.live for new victims and groups and
// persists them to the store, notifying callers of newly-observed victims.
package ingest

import (
	"context"

	"github.com/FLINTEK-LLC/omen/internal/store"
)

// Source fetches victim and group data from an upstream feed. The
// Ransomware.live client is the production implementation; tests may supply
// a fake.
type Source interface {
	FetchRecentVictims(ctx context.Context) ([]store.Victim, error)
	FetchGroups(ctx context.Context) ([]store.Group, error)
}
