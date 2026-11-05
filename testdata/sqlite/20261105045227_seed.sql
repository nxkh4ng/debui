-- +goose Up
INSERT INTO sql_users (id, name, email, score, is_active, created_at) VALUES
    (1, 'Nguyễn Văn A', 'a@example.com', 9.5, 1, '2026-01-02 08:30:00'),
    (2, 'O''Brien, "Bob"', 'bob@example.com', 7.25, 0, '2026-01-03 09:00:00'),
    (3, NULL, NULL, NULL, 1, '2026-01-04 10:15:00');

WITH RECURSIVE seq(n) AS (SELECT 4 UNION ALL SELECT n + 1 FROM seq WHERE n < 150)
INSERT INTO sql_users (id, name, email, score, is_active, created_at)
SELECT n, 'user ' || n, 'user' || n || '@example.com', n / 10.0, n % 2, '2026-02-01 00:00:00'
FROM seq;

INSERT INTO "sql_order" (id, total) VALUES (1, 19.99), (2, 0.0), (3, NULL);

INSERT INTO "sql mixed types"
    (id, small_int, big_int, price, ratio, title, note, flag, payload, created_at, updated_at, maybe_null)
VALUES
    (1, 42, 9007199254740993, 1234.5678, 0.000001, 'hello', 'line1
line2', 1, X'DEADBEEF', '2026-01-01T00:00:00Z', '2026-01-02 03:04:05', NULL),
    (2, -7, -9007199254740993, -0.5, NULL, '', NULL, 0, NULL, NULL, NULL, NULL);

UPDATE "sql_order" SET code = 'ORD-001' WHERE id = 1;
UPDATE "sql_order" SET code = 'ORD-002' WHERE id = 2;

-- +goose Down
DELETE FROM "sql mixed types";
DELETE FROM "sql_order";
DELETE FROM sql_users;
