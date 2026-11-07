package driver

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

type SQLite struct {
	db *sql.DB
}

var _ Driver = (*SQLite)(nil)

func OpenSQLite(ctx context.Context, dsn string) (*SQLite, error) {
	normDSN, err := normalizeDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("resolve db path %q: %w", dsn, err)
	}

	db, err := sql.Open("sqlite", normDSN)
	if err != nil {
		return nil, fmt.Errorf("open sqlite %q: %w", dsn, err)
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("open sqlite %q: %w", dsn, err)
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
		return nil, fmt.Errorf("list tables: %w", err)
	}
	defer func() { _ = rows.Close() }()

	tables := make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("scan tables: %w", err)
		}
		tables = append(tables, name)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tables: %w", err)
	}

	return tables, nil
}

func (s *SQLite) GetData(ctx context.Context, table string, limit, offset int) (TableData, error) {
	exist, err := s.tableExist(ctx, table)
	if err != nil {
		return TableData{}, err
	}
	if !exist {
		return TableData{}, ErrTableNotFound
	}
	q := quoteTable(table)

	var total int64
	if err := s.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM `+q,
	).Scan(&total); err != nil {
		return TableData{}, fmt.Errorf("count rows in %q: %w", table, err)
	}

	rows, err := s.db.QueryContext(ctx,
		`SELECT * FROM `+q+` LIMIT ? OFFSET ?`,
		limit, offset,
	)
	if err != nil {
		return TableData{}, fmt.Errorf("read table %q: %w", table, err)
	}
	defer func() { _ = rows.Close() }()

	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return TableData{}, fmt.Errorf("columns type int table %q: %w", table, err)
	}
	columns := make([]Column, len(colTypes))
	for i, ct := range colTypes {
		columns[i] = Column{Name: ct.Name(), Type: ct.DatabaseTypeName()}
	}

	vals := make([]any, len(columns))
	ptrs := make([]any, len(columns))
	for i := range vals {
		ptrs[i] = &vals[i]
	}

	rowsData := make([][]any, 0, limit)
	row := 0
	for rows.Next() {
		row++
		if err := rows.Scan(ptrs...); err != nil {
			return TableData{}, fmt.Errorf("scan talbe %q row %d: %w", table, row, err)
		}

		cells := make([]any, len(vals))
		for i, v := range vals {
			cells[i] = s.convertCell(v)
		}
		rowsData = append(rowsData, cells)
	}
	if err := rows.Err(); err != nil {
		return TableData{}, fmt.Errorf("iterate table %q: %w", table, err)
	}

	return TableData{
		Columns: columns,
		Rows:    rowsData,
		Total:   total,
	}, nil
}

func (s *SQLite) Close() error {
	return s.db.Close()
}

func (s *SQLite) tableExist(ctx context.Context, table string) (bool, error) {
	var exist int
	err := s.db.QueryRowContext(ctx,
		`SELECT 1
		FROM sqlite_master
		WHERE type IN ('table', 'view')
			AND name NOT GLOB 'sqlite_*'
			AND name = ?`,
		table,
	).Scan(&exist)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("check table %q: %w", table, err)
	default:
		return true, nil
	}
}

func (s *SQLite) convertCell(v any) any {
	switch x := v.(type) {
	case nil:
		return nil
	case int64:
		return strconv.FormatInt(x, 10)
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case string:
		return x
	case []byte:
		const maxConn = 64
		if len(x) <= maxConn {
			return "\\x" + hex.EncodeToString(x)
		}
		return "\\x" + hex.EncodeToString(x[:maxConn]) + fmt.Sprintf("… (%d bytes)", len(x))
	case time.Time:
		return x.Format(time.RFC3339Nano)
	default:
		return fmt.Sprintf("%v", x)
	}
}

func normalizeDSN(dsn string) (string, error) {
	abs, err := filepath.Abs(dsn)
	if err != nil {
		return "", err
	}

	u := &url.URL{
		Scheme:   "file",
		Path:     filepath.ToSlash(abs),
		RawQuery: "mode=rw&_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)",
	}

	return u.String(), nil
}
