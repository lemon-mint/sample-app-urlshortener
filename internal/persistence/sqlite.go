package persistence

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/user/urlshortener/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

const ( 
    base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    shortCodeLength = 6
)

// SQLitePersistence is a SQLite implementation of the URLPersistence interface.
type SQLitePersistence struct {
	db *sql.DB
}

// NewSQLitePersistence creates a new SQLitePersistence instance.
func NewSQLitePersistence(dataSourceName string) (*SQLitePersistence, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("persistence: failed to open database: %w", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		original TEXT NOT NULL UNIQUE,
		short TEXT NOT NULL UNIQUE
	);
	`)
	if err != nil {
		return nil, fmt.Errorf("persistence: failed to create table: %w", err)
	}

	return &SQLitePersistence{db: db}, nil
}

// Save saves a new URL record.
func (p *SQLitePersistence) Save(originalURL string) (string, error) {
    var shortCode string
    err := p.db.QueryRow("SELECT short FROM urls WHERE original = ?", originalURL).Scan(&shortCode)
    if err == nil {
        return shortCode, nil
    }

    if err != sql.ErrNoRows {
        return "", fmt.Errorf("persistence: failed to query url: %w", err)
    }

    shortCode = p.generateShortCode()

	_, err = p.db.Exec("INSERT INTO urls (original, short) VALUES (?, ?)", originalURL, shortCode)
	if err != nil {
		return "", fmt.Errorf("persistence: failed to insert url: %w", err)
	}

	return shortCode, nil
}

// Get retrieves the original URL for a given short code.
func (p *SQLitePersistence) Get(shortCode string) (string, error) {
	var originalURL string
	err := p.db.QueryRow("SELECT original FROM urls WHERE short = ?", shortCode).Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", types.ErrNotFound
		}
		return "", fmt.Errorf("persistence: failed to get url: %w", err)
	}

	return originalURL, nil
}

func (p *SQLitePersistence) generateShortCode() string {
    rand.Seed(time.Now().UnixNano())
    b := make([]byte, shortCodeLength)
    for i := range b {
        b[i] = base62Chars[rand.Intn(len(base62Chars))]
    }
    return string(b)
}
