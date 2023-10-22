
CREATE TABLE user_role (
    user_id UUID NOT NULL REFERENCES account(id) ON DELETE CASCADE,
    role VARCHAR(32) NOT NULL REFERENCES role_permission(role) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(user_id, role)
);

-- name: CreateUserRole :one
INSERT INTO user_role (user_id, role)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserRoleByUser :many
SELECT * FROM user_role WHERE user_id=$1;


-- name: DeleteUserRoleByUser :exec
DELETE FROM user_role WHERE user_id=$1;