-- name: CreateAccount :one
INSERT INTO account (id, email, password, first_name, last_name, is_administrator)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAccountByEmail :one
SELECT * FROM account
WHERE email=$1;

-- name: GetActiveAccountById :one
SELECT * FROM account
WHERE id=$1 and active=TRUE;