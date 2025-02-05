package types

import (
	"time"

	"github.com/lib/pq"
)

type Category struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Image       *string        `json:"image,omitempty"`
	Depth       int32          `json:"depth"`
	ParentID    *string        `json:"parent_id,omitempty"`
	Children    []*Category    `json:"children"`
	Path        pq.StringArray `json:"-"`
	CreatedAt   time.Time      `json:"created_at"`
}
