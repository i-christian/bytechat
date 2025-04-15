// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
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
returning user_id, last_name, first_name, gender, email, password, created_at, status, updated_at, role_id
`

type CreateUserParams struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
	Password  string `json:"password"`
	Name      string `json:"name"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.Gender,
		arg.Password,
		arg.Name,
	)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.LastName,
		&i.FirstName,
		&i.Gender,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.Status,
		&i.UpdatedAt,
		&i.RoleID,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
delete from users
where user_id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteUser, userID)
	return err
}

const editPassword = `-- name: EditPassword :exec
update users
    set password = coalesce($2, password)
where user_id = $1
`

type EditPasswordParams struct {
	UserID   uuid.UUID `json:"user_id"`
	Password string    `json:"password"`
}

func (q *Queries) EditPassword(ctx context.Context, arg EditPasswordParams) error {
	_, err := q.db.Exec(ctx, editPassword, arg.UserID, arg.Password)
	return err
}

const editUser = `-- name: EditUser :exec
update users
    set first_name = coalesce($2, first_name),
    last_name = coalesce($3, last_name),
    gender = coalesce($4, gender),
    email = coalesce($5, email)
where user_id = $1
`

type EditUserParams struct {
	UserID    uuid.UUID `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Gender    string    `json:"gender"`
	Email     string    `json:"email"`
}

func (q *Queries) EditUser(ctx context.Context, arg EditUserParams) error {
	_, err := q.db.Exec(ctx, editUser,
		arg.UserID,
		arg.FirstName,
		arg.LastName,
		arg.Gender,
		arg.Email,
	)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
select password, user_id from users 
where email = $1
`

type GetUserByEmailRow struct {
	Password string    `json:"password"`
	UserID   uuid.UUID `json:"user_id"`
}

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (GetUserByEmailRow, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i GetUserByEmailRow
	err := row.Scan(&i.Password, &i.UserID)
	return i, err
}

const getUserDetails = `-- name: GetUserDetails :one
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
    users.user_id = $1
`

type GetUserDetailsRow struct {
	UserID    uuid.UUID `json:"user_id"`
	LastName  string    `json:"last_name"`
	FirstName string    `json:"first_name"`
	Gender    string    `json:"gender"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
}

func (q *Queries) GetUserDetails(ctx context.Context, userID uuid.UUID) (GetUserDetailsRow, error) {
	row := q.db.QueryRow(ctx, getUserDetails, userID)
	var i GetUserDetailsRow
	err := row.Scan(
		&i.UserID,
		&i.LastName,
		&i.FirstName,
		&i.Gender,
		&i.Email,
		&i.Password,
		&i.Role,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
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
order by roles.name, users.last_name, users.first_name
`

type ListUsersRow struct {
	UserID    uuid.UUID `json:"user_id"`
	LastName  string    `json:"last_name"`
	FirstName string    `json:"first_name"`
	Gender    string    `json:"gender"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
}

func (q *Queries) ListUsers(ctx context.Context) ([]ListUsersRow, error) {
	rows, err := q.db.Query(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListUsersRow{}
	for rows.Next() {
		var i ListUsersRow
		if err := rows.Scan(
			&i.UserID,
			&i.LastName,
			&i.FirstName,
			&i.Gender,
			&i.Email,
			&i.Password,
			&i.Role,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const userOfflineStatus = `-- name: UserOfflineStatus :exec
update users
    set status = coalesce('offline', status)
where user_id = $1
`

func (q *Queries) UserOfflineStatus(ctx context.Context, userID uuid.UUID) error {
	_, err := q.db.Exec(ctx, userOfflineStatus, userID)
	return err
}

const userOnlineStatus = `-- name: UserOnlineStatus :exec
update users
    set status = coalesce('online', status)
where user_id = $1
`

func (q *Queries) UserOnlineStatus(ctx context.Context, userID uuid.UUID) error {
	_, err := q.db.Exec(ctx, userOnlineStatus, userID)
	return err
}
