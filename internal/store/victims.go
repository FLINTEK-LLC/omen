package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// VictimExists reports whether a victim with the given ID is already stored,
// which is how the poller distinguishes new events from ones already seen.
func (s *Store) VictimExists(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := s.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM victims WHERE id = ?)`, id,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check victim exists: %w", err)
	}
	return exists, nil
}

// UpsertVictim inserts a new victim or updates an existing one's mutable
// fields (press links, infostealer data, raw source) and last_updated stamp.
func (s *Store) UpsertVictim(ctx context.Context, v Victim) error {
	pressLinks, err := json.Marshal(v.PressLinks)
	if err != nil {
		return fmt.Errorf("marshal press_links: %w", err)
	}
	infostealer, err := json.Marshal(v.Infostealer)
	if err != nil {
		return fmt.Errorf("marshal infostealer: %w", err)
	}
	rawSource, err := json.Marshal(v.RawSource)
	if err != nil {
		return fmt.Errorf("marshal raw_source: %w", err)
	}

	now := time.Now().UTC()
	if v.FirstSeen.IsZero() {
		v.FirstSeen = now
	}

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO victims (
			id, victim_name, group_name, attack_date, country, domain,
			press_links, infostealer, raw_source, first_seen, last_updated
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			press_links  = excluded.press_links,
			infostealer  = excluded.infostealer,
			raw_source   = excluded.raw_source,
			domain       = excluded.domain,
			last_updated = excluded.last_updated
	`, v.ID, v.VictimName, v.GroupName, v.AttackDate, v.Country, v.Domain,
		string(pressLinks), string(infostealer), string(rawSource),
		v.FirstSeen, now,
	)
	if err != nil {
		return fmt.Errorf("upsert victim: %w", err)
	}
	return nil
}

// GetVictim returns a single victim by ID, or (Victim{}, sql.ErrNoRows) if
// not found.
func (s *Store) GetVictim(ctx context.Context, id string) (Victim, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, victim_name, group_name, attack_date, country, domain,
		       press_links, infostealer, raw_source, first_seen, last_updated
		FROM victims WHERE id = ?
	`, id)
	return scanVictim(row)
}

// VictimFilter narrows the results of ListVictims.
type VictimFilter struct {
	Country string
	Group   string
	Since   string // attack_date >= Since (YYYY-MM-DD)
	Limit   int
	Offset  int
}

// ListVictims returns victims matching filter, newest attack_date first.
func (s *Store) ListVictims(ctx context.Context, f VictimFilter) ([]Victim, error) {
	limit := f.Limit
	if limit <= 0 || limit > 500 {
		limit = 50
	}

	query := `
		SELECT id, victim_name, group_name, attack_date, country, domain,
		       press_links, infostealer, raw_source, first_seen, last_updated
		FROM victims WHERE 1=1
	`
	var args []any
	if f.Country != "" {
		query += " AND country = ?"
		args = append(args, f.Country)
	}
	if f.Group != "" {
		query += " AND group_name = ?"
		args = append(args, f.Group)
	}
	if f.Since != "" {
		query += " AND attack_date >= ?"
		args = append(args, f.Since)
	}
	query += " ORDER BY attack_date DESC, first_seen DESC LIMIT ? OFFSET ?"
	args = append(args, limit, f.Offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list victims: %w", err)
	}
	defer rows.Close()

	var out []Victim
	for rows.Next() {
		v, err := scanVictim(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanVictim(row rowScanner) (Victim, error) {
	var v Victim
	var pressLinks, infostealer, rawSource string
	err := row.Scan(
		&v.ID, &v.VictimName, &v.GroupName, &v.AttackDate, &v.Country, &v.Domain,
		&pressLinks, &infostealer, &rawSource, &v.FirstSeen, &v.LastUpdated,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Victim{}, err
		}
		return Victim{}, fmt.Errorf("scan victim: %w", err)
	}
	if err := json.Unmarshal([]byte(pressLinks), &v.PressLinks); err != nil {
		return Victim{}, fmt.Errorf("unmarshal press_links: %w", err)
	}
	if err := json.Unmarshal([]byte(infostealer), &v.Infostealer); err != nil {
		return Victim{}, fmt.Errorf("unmarshal infostealer: %w", err)
	}
	if err := json.Unmarshal([]byte(rawSource), &v.RawSource); err != nil {
		return Victim{}, fmt.Errorf("unmarshal raw_source: %w", err)
	}
	return v, nil
}
