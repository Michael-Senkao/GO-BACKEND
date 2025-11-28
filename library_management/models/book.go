package models

import "time"

// Book represents a library book.
type Book struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Status      string    `json:"status"` // "Available" or "Borrowed"
	ReservedBy  int       `json:"reserved_by,omitempty"`
	ReservedAt  time.Time `json:"reserved_at,omitempty"`
}
