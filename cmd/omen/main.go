// Command omen runs the OMEN ransomware CTI aggregation service: it polls
// Ransomware.live for new victims, enriches them with Shodan exposure data
// and CISA KEV status, fires watchlist alerts, and serves a REST API, SSE
// stream, and the live map dashboard.
package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/FLINTEK-LLC/omen/internal/api"
	"github.com/FLINTEK-LLC/omen/internal/config"
	"github.com/FLINTEK-LLC/omen/internal/enrich"
	"github.com/FLINTEK-LLC/omen/internal/ingest"
	"github.com/FLINTEK-LLC/omen/internal/store"
	"github.com/FLINTEK-LLC/omen/internal/watchlist"
	"github.com/FLINTEK-LLC/omen/web"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config.yaml (see config.sample.yaml)")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(cfg.Database.Path), 0o755); err != nil {
		log.Fatalf("create database dir: %v", err)
	}

	st, err := store.Open(cfg.Database.Path)
	if err != nil {
		log.Fatalf("open store: %v", err)
	}
	defer st.Close()

	kevCatalog := enrich.NewKEVCatalog(cfg.KEV.FeedURL, cfg.KEV.CachePath)
	if err := kevCatalog.LoadFromDisk(); err != nil {
		log.Printf("kev: load disk cache: %v", err)
	}

	var enricher enrich.Enricher
	if cfg.Shodan.APIKey == "" {
		log.Println("shodan: no API key configured, enrichment disabled")
		enricher = enrich.NoopEnricher{}
	} else {
		enricher = enrich.NewShodanClient(cfg.Shodan.APIKey, cfg.Shodan.RateLimitInterval.Duration)
	}

	notifier := watchlist.NewWebhookNotifier(cfg.Watchlist.DefaultNotifyVia)
	hub := api.NewHub()
	apiServer := api.NewServer(st, hub, web.Files)

	onNewVictim := func(ctx context.Context, v store.Victim) {
		if err := enrich.EnrichVictim(ctx, st, enricher, kevCatalog, v); err != nil {
			log.Printf("enrich victim %s: %v", v.ID, err)
		}

		entries, err := st.ListWatchlist(ctx)
		if err != nil {
			log.Printf("watchlist: list entries: %v", err)
		} else {
			for _, m := range watchlist.Match(entries, v.VictimName, v.Domain) {
				evt := watchlist.Event{
					VictimID:       v.ID,
					VictimName:     v.VictimName,
					GroupName:      v.GroupName,
					Domain:         v.Domain,
					MatchedPattern: m.Pattern,
					MatchedLabel:   m.Label,
				}
				if err := notifier.Notify(ctx, m.NotifyVia, evt); err != nil {
					log.Printf("watchlist: notify %q: %v", m.Pattern, err)
				}
			}
		}

		apiServer.BroadcastNewVictim(v)
	}

	source := ingest.NewRansomwareLiveClient(cfg.RansomwareLive.BaseURL)
	poller := ingest.NewPoller(source, st, cfg.Poll.VictimsInterval.Duration, cfg.Poll.GroupsInterval.Duration, onNewVictim)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go poller.Run(ctx)
	go runKEVRefreshLoop(ctx, kevCatalog, cfg.Poll.KEVInterval.Duration)

	httpServer := &http.Server{
		Addr:    cfg.Server.ListenAddr,
		Handler: apiServer.Handler(),
	}

	go func() {
		log.Printf("omen: listening on %s", cfg.Server.ListenAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("omen: shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("http server shutdown: %v", err)
	}
}

func runKEVRefreshLoop(ctx context.Context, kev *enrich.KEVCatalog, interval time.Duration) {
	refresh := func() {
		if err := kev.Refresh(ctx); err != nil {
			log.Printf("kev: refresh: %v", err)
		}
	}

	refresh()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			refresh()
		}
	}
}
