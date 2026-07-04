# OMEN

Ransomware CTI aggregation platform. Polls [Ransomware.live](https://www.ransomware.live)
for new victims and groups, enriches victims with Shodan exposure data,
cross-references exposed CVEs against the CISA Known Exploited
Vulnerabilities catalog, and serves a live attack map, dashboard, and
watchlist alerting system.

Part of the FLINTEK toolset (alongside CIRRUS, FIREAUDIT, and OBSERVER).

## Quick start (local)

```sh
cp config.sample.yaml config.yaml
go run ./cmd/omen -config config.yaml
```

Open http://localhost:8080 for the live map / dashboard.

By default Shodan enrichment is disabled (no-op) until you set an API key,
either in `config.yaml` under `shodan.api_key` or via the `SHODAN_API_KEY`
environment variable (preferred, since it keeps the key out of the config
file on disk).

## Docker

```sh
docker compose up --build
```

This runs with `config.sample.yaml`'s defaults (Shodan disabled, 30-minute
victim polling, SQLite + KEV cache under a named `omen-data` volume). To
supply a Shodan key:

```sh
SHODAN_API_KEY=yourkey docker compose up --build
```

To use a custom `config.yaml` instead of the bundled sample, uncomment the
config bind mount and `command` override in `docker-compose.yml`.

Building/running the image directly, without compose:

```sh
docker build -t omen .
docker run -p 8080:8080 -v omen-data:/app/data -e SHODAN_API_KEY=yourkey omen
```

## Configuration

See [config.sample.yaml](config.sample.yaml) for all options: server listen
address, database path, poll intervals, Ransomware.live/Shodan/CISA KEV
endpoints, and the default watchlist notification webhook.

## API

```
GET  /api/victims                 list, paginated, filter by date/country/group
GET  /api/victims/:id             full detail incl. shodan + kev data
GET  /api/groups                  list of tracked groups
GET  /api/groups/:name            group profile
GET  /api/watchlist               list watchlist entries
POST /api/watchlist                add entry
DEL  /api/watchlist/:id           remove entry
GET  /api/stream                  SSE stream of new-victim events
GET  /api/stats                   summary counts for dashboard (by country/group/day)
```

## Repo structure

```
omen/
├── cmd/omen/            # main entrypoint
├── internal/
│   ├── ingest/          # Ransomware.live poller
│   ├── enrich/          # Shodan + KEV logic
│   ├── watchlist/       # alerting
│   ├── store/           # SQLite access layer
│   ├── config/          # YAML config loading
│   └── api/             # REST + SSE handlers
├── web/                 # static frontend, embedded into the binary
├── config.sample.yaml
├── Dockerfile
├── docker-compose.yml
├── LICENSE (Apache 2.0)
└── README.md
```

## Known v1 limitations

- Domain resolution for victims without a clean domain from Ransomware.live
  is not yet implemented -- `Domain` is left blank, which means Shodan
  enrichment is skipped for those victims. A name-to-domain guesser or manual
  override is the planned follow-up.
- Watchlist matching is a deliberately dumb case-insensitive substring match
  (no fuzzy matching library), per the v1 spec.
- No auth/RBAC on the API -- run behind your own reverse proxy/VPN if
  exposing beyond a trusted network.

## License

Apache 2.0, see [LICENSE](LICENSE). Ransomware.live data is cached locally
for enrichment purposes; confirm their terms before redistributing raw
payloads long-term.
