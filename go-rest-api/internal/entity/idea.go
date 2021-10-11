package entity

import (
	"github.com/lib/pq"
	"time"
)

// Idea represents an idea record.
type Idea struct {
	ID           string    `json:"id"`
	AuthorEmail string    `json:"author_email"`
	Tags         pq.StringArray    `json:"tags"`
	Summary      string    `json:"summary"`
	Content      string    `json:"content"`
	Media        pq.StringArray    `json:"media"`
	MediaTypes   pq.StringArray    `json:"media_types"`
	BadFlag      bool    `json:"bad_flag"`
	Enabled      bool    `json:"enabled"`
	Issues       pq.StringArray   `json:"issues"`
	IssuesIPs    pq.StringArray    `json:"issues_ips"`
	Votes        int    `json:"votes"`
	VotersIds    pq.StringArray    `json:"voters_ids"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}