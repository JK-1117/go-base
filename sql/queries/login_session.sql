-- name: CreateLoginSession :one
INSERT INTO login_session (session_id, user_id, last_login, ip_addr, user_agent, expired_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetSessionBySessionId :one
SELECT * FROM login_session WHERE session_id=$1;


-- name: DeleteSessionBySessionId :exec
DELETE FROM login_session WHERE session_id=$1;


-- name: DeleteExpiredSession :exec
DELETE FROM login_session WHERE expired_at < now();