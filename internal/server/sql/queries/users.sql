-- name: CreateUser :one
insert into users (first_name, last_name, email, gender, password, role_id)
values (
    $1,
    $2,
    $3,
    $4,
    $5,
    (select role_id from roles where name = $6)
)
on conflict (email) do nothing
returning *;

-- name: UserOnlineStatus :exec
update users
    set status = coalesce('online', status)
where user_id = $1;

-- name: UserOfflineStatus :exec
update users
    set status = coalesce('offline', status)
where user_id = $1;

-- name: GetUserDetails :one
select
    users.user_id,
    users.last_name, 
    users.first_name, 
    users.gender, 
    users.email, 
    users.password, 
    roles.name as role
from 
    users
inner join 
    roles 
on 
    users.role_id = roles.role_id
where 
    users.user_id = $1;

-- name: GetUserByEmail :one
SELECT password, user_id FROM users 
WHERE email = $1;

-- name: ListUsers :many
SELECT
    users.user_id,
    users.last_name,
    users.first_name,
    users.gender,
    users.email,
    users.password,
    roles.name AS role
FROM users
INNER JOIN roles ON users.role_id = roles.role_id
ORDER BY last_name;

-- name: EditUser :exec
UPDATE users
    set first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    gender = COALESCE($4, gender),
    email = COALESCE($5, email)
WHERE user_id = $1;

-- name: EditPassword :exec
UPDATE users
    set password = COALESCE($2, password)
WHERE user_id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1;
