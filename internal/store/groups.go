package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// UpsertGroup inserts or refreshes a cached group profile.
func (s *Store) UpsertGroup(ctx context.Context, g Group) error {
	profile, err := json.Marshal(g.ProfileJSON)
	if err != nil {
		return fmt.Errorf("marshal profile_json: %w", err)
	}
	now := time.Now().UTC()
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO groups (name, profile_json, last_updated)
		VALUES (?, ?, ?)
		ON CONFLICT(name) DO UPDATE SET
			profile_json = excluded.profile_json,
			last_updated = excluded.last_updated
	`, g.Name, string(profile), now)
	if err != nil {
		return fmt.Errorf("upsert group: %w", err)
	}
	return nil
}

// GetGroup returns a single cached group profile by name.
func (s *Store) GetGroup(ctx context.Context, name string) (Group, error) {
	var g Group
	var profile string
	err := s.db.QueryRowContext(ctx, `
		SELECT name, profile_json, last_updated FROM groups WHERE name = ?
	`, name).Scan(&g.Name, &profile, &g.LastUpdated)
	if err != nil {
		return Group{}, err
	}
	if err := json.Unmarshal([]byte(profile), &g.ProfileJSON); err != nil {
		return Group{}, fmt.Errorf("unmarshal profile_json: %w", err)
	}
	return g, nil
}

// ListGroups returns all cached group profiles.
func (s *Store) ListGroups(ctx context.Context) ([]Group, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT name, profile_json, last_updated FROM groups ORDER BY name
	`)
	if err != nil {
		return nil, fmt.Errorf("list groups: %w", err)
	}
	defer rows.Close()

	var out []Group
	for rows.Next() {
		var g Group
		var profile string
		if err := rows.Scan(&g.Name, &profile, &g.LastUpdated); err != nil {
			return nil, fmt.Errorf("scan group: %w", err)
		}
		if err := json.Unmarshal([]byte(profile), &g.ProfileJSON); err != nil {
			return nil, fmt.Errorf("unmarshal profile_json: %w", err)
		}
		out = append(out, g)
	}
	return out, rows.Err()
}
