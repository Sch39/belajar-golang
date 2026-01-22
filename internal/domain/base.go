// internal/domain/base.go
package domain

import "time"

type Timestamp struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SoftDelete struct {
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	IsActive  bool  `json:"is_active"`
}