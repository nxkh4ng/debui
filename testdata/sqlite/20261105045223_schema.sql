-- +goose Up
CREATE TABLE sql_users (
    id         INTEGER PRIMARY KEY,
    name       TEXT,
    email      TEXT,
    score      REAL,
    is_active  INTEGER NOT NULL DEFAULT 1,
    created_at TEXT    NOT NULL
);

CREATE TABLE sql_empty_table (
    id   INTEGER PRIMARY KEY,
    note TEXT
);

CREATE TABLE "sql_order" (
    id    INTEGER PRIMARY KEY,
    total REAL,
    code  TEXT UNIQUE
);

CREATE TABLE "sql mixed types" (
    id          INTEGER PRIMARY KEY,
    small_int   INTEGER,
    big_int     INTEGER,
    price       REAL,
    ratio       REAL,
    title       TEXT,
    note        TEXT,
    flag        INTEGER,
    payload     BLOB,
    created_at  TEXT,
    updated_at  TIMESTAMP,
    maybe_null  TEXT
);

-- Index đơn giản, 1 cột
CREATE INDEX idx_users_created_at ON sql_users (created_at);

-- Index nhiều cột, sort ASC/DESC
CREATE INDEX idx_users_score_active ON sql_users (score ASC, is_active DESC);

-- Index unique (trả về unique = 1 trong pragma index_list)
CREATE UNIQUE INDEX idx_users_email ON sql_users (email);

-- Partial index (WHERE clause)
CREATE INDEX idx_users_active_only ON sql_users (name) WHERE is_active = 1;

-- Expression index (cột là biểu thức, không phải cột thuần)
CREATE INDEX idx_users_name_lower ON sql_users (lower(name));

-- Index với COLLATE
CREATE INDEX idx_users_email_nocase ON sql_users (email COLLATE NOCASE);

-- Index trên bảng có tên chứa khoảng trắng / cần quote
CREATE INDEX idx_mixed_title ON "sql mixed types" (title);

-- Index multi-column trên bảng quoted, có expression
CREATE INDEX idx_mixed_price_round ON "sql mixed types" (small_int, round(price));

-- Index trên bảng rỗng (vẫn phải liệt kê được)
CREATE INDEX idx_empty_note ON sql_empty_table (note);

CREATE VIEW sql_active_users AS
SELECT id, name, email FROM sql_users WHERE is_active = 1;

-- +goose Down
DROP INDEX IF EXISTS idx_empty_note;
DROP INDEX IF EXISTS idx_mixed_price_round;
DROP INDEX IF EXISTS idx_mixed_title;
DROP INDEX IF EXISTS idx_users_email_nocase;
DROP INDEX IF EXISTS idx_users_name_lower;
DROP INDEX IF EXISTS idx_users_active_only;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_score_active;
DROP INDEX IF EXISTS idx_users_created_at;

DROP VIEW sql_active_users;
DROP TABLE "sql mixed types";
DROP TABLE "sql_order";
DROP TABLE sql_empty_table;
DROP TABLE sql_users;
