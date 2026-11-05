package driver

import (
	"context"
	"database/sql"
	"net/url"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type SQLite struct {
	db *sql.DB
}

var _ Driver = (*SQLite)(nil)

func OpenSQLite(ctx context.Context, dsn string) (*SQLite, error) {
	normDSN, err := normalizeDSN(dsn)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", normDSN)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &SQLite{db: db}, nil
}

func (s *SQLite) ListTables(ctx context.Context) ([]string, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT name FROM sqlite_master
		WHERE type = 'table'
		AND name NOT GLOB 'sqlite_*'
		ORDER BY name`,
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	tables := make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		tables = append(tables, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

func (s *SQLite) Close() error {
	return s.db.Close()
}

func normalizeDSN(dsn string) (string, error) {
	abs, err := filepath.Abs(dsn)
	if err != nil {
		return "", err
	}

	u := &url.URL{
		Scheme:   "file",
		Path:     filepath.ToSlash(abs),
		RawQuery: "mode=rw",
	}

	return u.String(), nil
}
