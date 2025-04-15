// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: settings.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const editMyPassword = `-- name: EditMyPassword :exec
UPDATE users
    set password = COALESCE($2, password)
WHERE user_id = $1
`

type EditMyPasswordParams struct {
	UserID   uuid.UUID `json:"user_id"`
	Password string    `json:"password"`
}

func (q *Queries) EditMyPassword(ctx context.Context, arg EditMyPasswordParams) error {
	_, err := q.db.Exec(ctx, editMyPassword, arg.UserID, arg.Password)
	return err
}

const editMyProfile = `-- name: EditMyProfile :exec
UPDATE users
    set first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    gender = COALESCE($4, gender),
    email = COALESCE($5, email)
WHERE user_id = $1
`

type EditMyProfileParams struct {
	UserID    uuid.UUID `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Gender    string    `json:"gender"`
	Email     string    `json:"email"`
}

func (q *Queries) EditMyProfile(ctx context.Context, arg EditMyProfileParams) error {
	_, err := q.db.Exec(ctx, editMyProfile,
		arg.UserID,
		arg.FirstName,
		arg.LastName,
		arg.Gender,
		arg.Email,
	)
	return err
}
