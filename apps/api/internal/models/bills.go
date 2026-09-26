package models

import "time"

// Bill represents a single bill record.
type Bill struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Cost        float64   `json:"cost"`
	DatePaid    time.Time `json:"date_paid"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
