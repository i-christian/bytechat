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
select password, user_id from users 
where email = $1;

-- name: ListUsers :many
select
    users.user_id,
    users.last_name,
    users.first_name,
    users.gender,
    users.email,
    users.password,
    roles.name as role
from users
join roles on users.role_id = roles.role_id
order by roles.name, users.last_name, users.first_name;

-- name: EditUser :exec
update users
    set first_name = coalesce($2, first_name),
    last_name = coalesce($3, last_name),
    gender = coalesce($4, gender),
    email = coalesce($5, email)
where user_id = $1;

-- name: EditPassword :exec
update users
    set password = coalesce($2, password)
where user_id = $1;

-- name: DeleteUser :exec
delete from users
where user_id = $1;
