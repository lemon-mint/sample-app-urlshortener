package types

import "errors"

var (
	// ErrNotFound is returned when a requested record is not found.
	ErrNotFound = errors.New("types: record not found")
)

// URL represents a shortened URL record.
type URL struct {
	ID       int64  `json:"id"`
	Original string `json:"original"`
	Short    string `json:"short"`
}

// URLPersistence defines the contract for URL persistence operations.
type URLPersistence interface {
	// Save saves a new URL record.
	Save(originalURL string) (string, error)
	// Get retrieves the original URL for a given short code.
	Get(shortCode string) (string, error)
}
