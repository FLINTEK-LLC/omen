package store

import (
	"context"
	"fmt"
)

// Stats is the summary payload backing the dashboard.
type Stats struct {
	TotalVictims int            `json:"total_victims"`
	ByCountry    map[string]int `json:"by_country"`
	ByGroup      map[string]int `json:"by_group"`
	ByDay        map[string]int `json:"by_day"`
}

// GetStats computes summary counts across all stored victims.
func (s *Store) GetStats(ctx context.Context) (Stats, error) {
	stats := Stats{
		ByCountry: map[string]int{},
		ByGroup:   map[string]int{},
		ByDay:     map[string]int{},
	}

	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM victims`).Scan(&stats.TotalVictims); err != nil {
		return Stats{}, fmt.Errorf("count victims: %w", err)
	}

	if err := scanCounts(ctx, s, `
		SELECT COALESCE(NULLIF(country, ''), 'unknown'), COUNT(*)
		FROM victims GROUP BY 1
	`, stats.ByCountry); err != nil {
		return Stats{}, fmt.Errorf("count by country: %w", err)
	}

	if err := scanCounts(ctx, s, `
		SELECT group_name, COUNT(*) FROM victims GROUP BY group_name
	`, stats.ByGroup); err != nil {
		return Stats{}, fmt.Errorf("count by group: %w", err)
	}

	if err := scanCounts(ctx, s, `
		SELECT COALESCE(NULLIF(attack_date, ''), 'unknown'), COUNT(*)
		FROM victims GROUP BY 1
	`, stats.ByDay); err != nil {
		return Stats{}, fmt.Errorf("count by day: %w", err)
	}

	return stats, nil
}

func scanCounts(ctx context.Context, s *Store, query string, dest map[string]int) error {
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		var count int
		if err := rows.Scan(&key, &count); err != nil {
			return err
		}
		dest[key] = count
	}
	return rows.Err()
}
