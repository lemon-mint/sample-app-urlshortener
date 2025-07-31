package core

import (
	"fmt"

	"github.com/user/urlshortener/internal/types"
)

// Core is the core business logic of the URL shortener.
type Core struct {
	persistence types.URLPersistence
}

// NewCore creates a new Core instance.
func NewCore(persistence types.URLPersistence) *Core {
	return &Core{persistence: persistence}
}

// ShortenURL shortens a URL.
func (c *Core) ShortenURL(originalURL string) (string, error) {
	shortCode, err := c.persistence.Save(originalURL)
	if err != nil {
		return "", fmt.Errorf("core: failed to shorten url: %w", err)
	}
	return shortCode, nil
}

// GetURL retrieves the original URL for a given short code.
func (c *Core) GetURL(shortCode string) (string, error) {
	originalURL, err := c.persistence.Get(shortCode)
	if err != nil {
		return "", fmt.Errorf("core: failed to get url: %w", err)
	}
	return originalURL, nil
}
