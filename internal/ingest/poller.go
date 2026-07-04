package ingest

import (
	"context"
	"log"
	"time"

	"github.com/FLINTEK-LLC/omen/internal/store"
)

// NewVictimHandler is invoked for each victim the poller observes for the
// first time, after it has been persisted. Callers use this hook to trigger
// enrichment, watchlist matching, and SSE broadcast without the ingest
// package needing to know about them.
type NewVictimHandler func(ctx context.Context, v store.Victim)

// Poller periodically pulls victims and groups from a Source and persists
// them, diffing against the store to detect new victims.
type Poller struct {
	source          Source
	store           *store.Store
	victimsInterval time.Duration
	groupsInterval  time.Duration
	onNewVictim     NewVictimHandler
}

// NewPoller builds a Poller. onNewVictim may be nil if the caller doesn't
// need new-victim notifications.
func NewPoller(source Source, st *store.Store, victimsInterval, groupsInterval time.Duration, onNewVictim NewVictimHandler) *Poller {
	return &Poller{
		source:          source,
		store:           st,
		victimsInterval: victimsInterval,
		groupsInterval:  groupsInterval,
		onNewVictim:     onNewVictim,
	}
}

// Run starts the victim and group polling loops. It blocks until ctx is
// canceled.
func (p *Poller) Run(ctx context.Context) {
	go p.loop(ctx, p.victimsInterval, p.pollVictimsOnce)
	p.loop(ctx, p.groupsInterval, p.pollGroupsOnce)
}

func (p *Poller) loop(ctx context.Context, interval time.Duration, poll func(context.Context)) {
	poll(ctx)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			poll(ctx)
		}
	}
}

func (p *Poller) pollVictimsOnce(ctx context.Context) {
	victims, err := p.source.FetchRecentVictims(ctx)
	if err != nil {
		log.Printf("ingest: fetch recent victims: %v", err)
		return
	}

	for _, v := range victims {
		exists, err := p.store.VictimExists(ctx, v.ID)
		if err != nil {
			log.Printf("ingest: check victim %s exists: %v", v.ID, err)
			continue
		}
		if err := p.store.UpsertVictim(ctx, v); err != nil {
			log.Printf("ingest: upsert victim %s: %v", v.ID, err)
			continue
		}
		if !exists && p.onNewVictim != nil {
			p.onNewVictim(ctx, v)
		}
	}
}

func (p *Poller) pollGroupsOnce(ctx context.Context) {
	groups, err := p.source.FetchGroups(ctx)
	if err != nil {
		log.Printf("ingest: fetch groups: %v", err)
		return
	}

	for _, g := range groups {
		if err := p.store.UpsertGroup(ctx, g); err != nil {
			log.Printf("ingest: upsert group %s: %v", g.Name, err)
		}
	}
}
