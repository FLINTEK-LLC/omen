package store

import (
	"context"
	"fmt"
	"time"
)

// InsertWatchlistEntry adds a new watchlist pattern and returns its ID.
func (s *Store) InsertWatchlistEntry(ctx context.Context, e WatchlistEntry) (int64, error) {
	now := time.Now().UTC()
	res, err := s.db.ExecContext(ctx, `
		INSERT INTO watchlist (pattern, label, notify_via, created_at)
		VALUES (?, ?, ?, ?)
	`, e.Pattern, e.Label, e.NotifyVia, now)
	if err != nil {
		return 0, fmt.Errorf("insert watchlist entry: %w", err)
	}
	return res.LastInsertId()
}

// ListWatchlist returns all configured watchlist entries.
func (s *Store) ListWatchlist(ctx context.Context) ([]WatchlistEntry, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, pattern, label, notify_via, created_at FROM watchlist ORDER BY id
	`)
	if err != nil {
		return nil, fmt.Errorf("list watchlist: %w", err)
	}
	defer rows.Close()

	var out []WatchlistEntry
	for rows.Next() {
		var e WatchlistEntry
		if err := rows.Scan(&e.ID, &e.Pattern, &e.Label, &e.NotifyVia, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan watchlist entry: %w", err)
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// DeleteWatchlistEntry removes a watchlist entry by ID.
func (s *Store) DeleteWatchlistEntry(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM watchlist WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete watchlist entry: %w", err)
	}
	return nil
}
