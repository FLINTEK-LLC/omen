package store

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"
)

// Victim is a single ransomware victim record.
type Victim struct {
	ID           string
	VictimName   string
	GroupName    string
	AttackDate   string
	Country      string
	Domain       string
	PressLinks   []string
	Infostealer  map[string]any
	RawSource    map[string]any
	FirstSeen    time.Time
	LastUpdated  time.Time
}

// VictimID derives the stable primary key for a victim from its identifying
// fields. The same (victim, group, date) tuple always yields the same ID,
// which is what the poller uses to diff new records against existing ones.
func VictimID(victimName, groupName, attackDate string) string {
	key := strings.ToLower(strings.TrimSpace(victimName)) + "|" +
		strings.ToLower(strings.TrimSpace(groupName)) + "|" +
		strings.TrimSpace(attackDate)
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])[:24]
}

// ShodanSnapshot is a point-in-time Shodan lookup result for a victim's
// resolved IP/domain.
type ShodanSnapshot struct {
	ID        int64
	VictimID  string
	IP        string
	Ports     []int
	CVEs      []string
	QueriedAt time.Time
}

// KEVMatch records that one of a victim's Shodan-flagged CVEs is present in
// the CISA Known Exploited Vulnerabilities catalog.
type KEVMatch struct {
	ID        int64
	VictimID  string
	CVEID     string
	KEVAdded  string
	MatchedAt time.Time
}

// WatchlistEntry is a user-configured pattern to alert on.
type WatchlistEntry struct {
	ID        int64
	Pattern   string
	Label     string
	NotifyVia string
	CreatedAt time.Time
}

// Group is a cached ransomware group profile.
type Group struct {
	Name        string
	ProfileJSON map[string]any
	LastUpdated time.Time
}
