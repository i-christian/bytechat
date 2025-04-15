-- name: EditMyProfile :exec
UPDATE users
    set first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    gender = COALESCE($4, gender),
    email = COALESCE($5, email)
WHERE user_id = $1;

-- name: EditMyPassword :exec
UPDATE users
    set password = COALESCE($2, password)
WHERE user_id = $1;

