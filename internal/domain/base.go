// internal/domain/base.go
package domain

import "time"

type Timestamp struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SoftDelete struct {
	DeletedAt *time.Time
	IsActive  bool
}
