package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"read-adviser/internal/lib/e"
	"read-adviser/internal/storage"
)

type Storage struct {
	db *sql.DB
}

// New creates a new Sqlite storage
func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect database %w", err)
	}

	return &Storage{db: db}, nil
}

// Save saves page to storage.
func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	q := `INSERT INTO pages (url, user_name) VALUES (?, ?)`

	if _, err := s.db.ExecContext(ctx, q, p.URL, p.UserName); err != nil {
		return e.WrapIfErr("can't save page", err)
	}

	return nil
}

// PickRandom pick random page from storage.
func (s *Storage) PickRandom(ctx context.Context, userName string) (*storage.Page, error) {
	q := `SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM() LIMIT 1`

	var url string

	err := s.db.QueryRowContext(ctx, q, userName).Scan(&url)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, storage.ErrNoSavedPages
	}
	if err != nil {
		return nil, e.WrapIfErr("can't pick random page", err)
	}

	return &storage.Page{
		URL:      url,
		UserName: userName,
	}, nil
}

// Remove removes page from storage.
func (s *Storage) Remove(ctx context.Context, page *storage.Page) error {
	q := `DELETE FROM pages WHERE url = ? AND user_name = ?`
	if _, err := s.db.ExecContext(ctx, q, page.URL, page.UserName); err != nil {
		return fmt.Errorf("can't remove page: %w", err)
	}

	return nil
}

// IsExists checks if page exists in storage.
func (s *Storage) IsExists(ctx context.Context, page *storage.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE url =? AND user_name =?`

	var count int64

	err := s.db.QueryRowContext(ctx, q, page.URL, page.UserName).Scan(&count)
	if err != nil {
		return false, e.WrapIfErr("can't check if page exist", err)
	}

	return count > 0, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        url TEXT NOT NULL,
        user_name TEXT NOT NULL
    )`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return e.WrapIfErr("can't init database", err)
	}

	return nil
}
