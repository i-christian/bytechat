-- name: CreateSession :exec
INSERT INTO sessions (session_id, user_id) 
VALUES ($1, $2)
ON CONFLICT (user_id) 
DO UPDATE SET session_id = EXCLUDED.session_id;

-- name: GetSession :one
SELECT 
  sessions.user_id,
  sessions.session_id,
  roles.name AS role,
  sessions.expires
FROM sessions
INNER JOIN users
  ON sessions.user_id = users.user_id
INNER JOIN roles 
  ON users.role_id = roles.role_id
WHERE session_id = $1;

-- name: RefreshSession :exec
UPDATE sessions
  SET expires = COALESCE($2, expires),
  session_id = COALESCE($3, session_id)
WHERE user_id = $1;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE user_id = $1;
