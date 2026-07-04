CREATE TABLE IF NOT EXISTS victims (
    id            TEXT PRIMARY KEY,
    victim_name   TEXT NOT NULL,
    group_name    TEXT NOT NULL,
    attack_date   TEXT,
    country       TEXT,
    domain        TEXT,
    press_links   TEXT NOT NULL DEFAULT '[]',
    infostealer   TEXT NOT NULL DEFAULT '{}',
    raw_source    TEXT NOT NULL DEFAULT '{}',
    first_seen    DATETIME NOT NULL,
    last_updated  DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_victims_group ON victims(group_name);
CREATE INDEX IF NOT EXISTS idx_victims_country ON victims(country);
CREATE INDEX IF NOT EXISTS idx_victims_attack_date ON victims(attack_date);

CREATE TABLE IF NOT EXISTS shodan_snapshots (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    victim_id  TEXT NOT NULL REFERENCES victims(id) ON DELETE CASCADE,
    ip         TEXT NOT NULL,
    ports      TEXT NOT NULL DEFAULT '[]',
    cves       TEXT NOT NULL DEFAULT '[]',
    queried_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_shodan_victim ON shodan_snapshots(victim_id);

CREATE TABLE IF NOT EXISTS kev_matches (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    victim_id  TEXT NOT NULL REFERENCES victims(id) ON DELETE CASCADE,
    cve_id     TEXT NOT NULL,
    kev_added  TEXT,
    matched_at DATETIME NOT NULL,
    UNIQUE(victim_id, cve_id)
);

CREATE INDEX IF NOT EXISTS idx_kev_victim ON kev_matches(victim_id);

CREATE TABLE IF NOT EXISTS watchlist (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    pattern    TEXT NOT NULL,
    label      TEXT NOT NULL DEFAULT '',
    notify_via TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS groups (
    name          TEXT PRIMARY KEY,
    profile_json  TEXT NOT NULL DEFAULT '{}',
    last_updated  DATETIME NOT NULL
);
