package store

import (
	"context"
	"encoding/json"
	"fmt"
)

// InsertShodanSnapshot records a Shodan lookup result for a victim.
func (s *Store) InsertShodanSnapshot(ctx context.Context, snap ShodanSnapshot) error {
	ports, err := json.Marshal(snap.Ports)
	if err != nil {
		return fmt.Errorf("marshal ports: %w", err)
	}
	cves, err := json.Marshal(snap.CVEs)
	if err != nil {
		return fmt.Errorf("marshal cves: %w", err)
	}
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO shodan_snapshots (victim_id, ip, ports, cves, queried_at)
		VALUES (?, ?, ?, ?, ?)
	`, snap.VictimID, snap.IP, string(ports), string(cves), snap.QueriedAt)
	if err != nil {
		return fmt.Errorf("insert shodan snapshot: %w", err)
	}
	return nil
}

// ShodanSnapshotsForVictim returns all Shodan snapshots recorded for a victim,
// most recent first.
func (s *Store) ShodanSnapshotsForVictim(ctx context.Context, victimID string) ([]ShodanSnapshot, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, victim_id, ip, ports, cves, queried_at
		FROM shodan_snapshots WHERE victim_id = ? ORDER BY queried_at DESC
	`, victimID)
	if err != nil {
		return nil, fmt.Errorf("query shodan snapshots: %w", err)
	}
	defer rows.Close()

	var out []ShodanSnapshot
	for rows.Next() {
		var snap ShodanSnapshot
		var ports, cves string
		if err := rows.Scan(&snap.ID, &snap.VictimID, &snap.IP, &ports, &cves, &snap.QueriedAt); err != nil {
			return nil, fmt.Errorf("scan shodan snapshot: %w", err)
		}
		if err := json.Unmarshal([]byte(ports), &snap.Ports); err != nil {
			return nil, fmt.Errorf("unmarshal ports: %w", err)
		}
		if err := json.Unmarshal([]byte(cves), &snap.CVEs); err != nil {
			return nil, fmt.Errorf("unmarshal cves: %w", err)
		}
		out = append(out, snap)
	}
	return out, rows.Err()
}

// InsertKEVMatch records that a CVE found on a victim's exposed surface is a
// known-exploited vulnerability. Re-matching the same (victim, cve) pair is a
// no-op.
func (s *Store) InsertKEVMatch(ctx context.Context, m KEVMatch) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO kev_matches (victim_id, cve_id, kev_added, matched_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(victim_id, cve_id) DO NOTHING
	`, m.VictimID, m.CVEID, m.KEVAdded, m.MatchedAt)
	if err != nil {
		return fmt.Errorf("insert kev match: %w", err)
	}
	return nil
}

// KEVMatchesForVictim returns all known-exploited-vulnerability matches for a
// victim.
func (s *Store) KEVMatchesForVictim(ctx context.Context, victimID string) ([]KEVMatch, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, victim_id, cve_id, kev_added, matched_at
		FROM kev_matches WHERE victim_id = ? ORDER BY matched_at DESC
	`, victimID)
	if err != nil {
		return nil, fmt.Errorf("query kev matches: %w", err)
	}
	defer rows.Close()

	var out []KEVMatch
	for rows.Next() {
		var m KEVMatch
		if err := rows.Scan(&m.ID, &m.VictimID, &m.CVEID, &m.KEVAdded, &m.MatchedAt); err != nil {
			return nil, fmt.Errorf("scan kev match: %w", err)
		}
		out = append(out, m)
	}
	return out, rows.Err()
}
