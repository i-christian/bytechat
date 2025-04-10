// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: rooms.sql

package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createPublicRoom = `-- name: CreatePublicRoom :exec
insert into rooms(name, description, room_type) 
values ($1, $2, $3)
`

type CreatePublicRoomParams struct {
	Name        string      `json:"name"`
	Description pgtype.Text `json:"description"`
	RoomType    string      `json:"room_type"`
}

func (q *Queries) CreatePublicRoom(ctx context.Context, arg CreatePublicRoomParams) error {
	_, err := q.db.Exec(ctx, createPublicRoom, arg.Name, arg.Description, arg.RoomType)
	return err
}

const getPrivateRooms = `-- name: GetPrivateRooms :many
select users.first_name || ' ' || users.last_name as my_name,
    rooms.room_id
from users 
join chat_rooms on users.user_id = chat_rooms.user_id
join rooms on chat_rooms.room_id = rooms.room_id
where users.user_id = $1
and rooms.room_type = 'private'
order by rooms.room_id
`

type GetPrivateRoomsRow struct {
	MyName interface{} `json:"my_name"`
	RoomID uuid.UUID   `json:"room_id"`
}

func (q *Queries) GetPrivateRooms(ctx context.Context, userID uuid.UUID) ([]GetPrivateRoomsRow, error) {
	rows, err := q.db.Query(ctx, getPrivateRooms, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPrivateRoomsRow{}
	for rows.Next() {
		var i GetPrivateRoomsRow
		if err := rows.Scan(&i.MyName, &i.RoomID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRoomDetails = `-- name: GetRoomDetails :one
select name from rooms
where room_id = $1
and room_type = 'public'
`

func (q *Queries) GetRoomDetails(ctx context.Context, roomID uuid.UUID) (string, error) {
	row := q.db.QueryRow(ctx, getRoomDetails, roomID)
	var name string
	err := row.Scan(&name)
	return name, err
}

const getUsersInRoom = `-- name: GetUsersInRoom :many
select 
    users.first_name || ' ' || users.last_name as full_name,
    users.status,
    users.user_id,
    users.updated_at,
    rooms.name
from users
join chat_rooms using(user_id)
join rooms using (room_id)
where rooms.room_id = $1
order by users.updated_at desc
`

type GetUsersInRoomRow struct {
	FullName  interface{}        `json:"full_name"`
	Status    string             `json:"status"`
	UserID    uuid.UUID          `json:"user_id"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
	Name      string             `json:"name"`
}

func (q *Queries) GetUsersInRoom(ctx context.Context, roomID uuid.UUID) ([]GetUsersInRoomRow, error) {
	rows, err := q.db.Query(ctx, getUsersInRoom, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetUsersInRoomRow{}
	for rows.Next() {
		var i GetUsersInRoomRow
		if err := rows.Scan(
			&i.FullName,
			&i.Status,
			&i.UserID,
			&i.UpdatedAt,
			&i.Name,
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

const initiatePrivateRoom = `-- name: InitiatePrivateRoom :one
select from create_private_room($1, $2)
`

type InitiatePrivateRoomParams struct {
	UserA uuid.UUID `json:"user_a"`
	UserB uuid.UUID `json:"user_b"`
}

type InitiatePrivateRoomRow struct {
}

func (q *Queries) InitiatePrivateRoom(ctx context.Context, arg InitiatePrivateRoomParams) (InitiatePrivateRoomRow, error) {
	row := q.db.QueryRow(ctx, initiatePrivateRoom, arg.UserA, arg.UserB)
	var i InitiatePrivateRoomRow
	err := row.Scan()
	return i, err
}

const joinRoom = `-- name: JoinRoom :exec
insert into chat_rooms(user_id, room_id) 
values ($1, $2)
on conflict (user_id, room_id) do nothing
`

type JoinRoomParams struct {
	UserID uuid.UUID `json:"user_id"`
	RoomID uuid.UUID `json:"room_id"`
}

func (q *Queries) JoinRoom(ctx context.Context, arg JoinRoomParams) error {
	_, err := q.db.Exec(ctx, joinRoom, arg.UserID, arg.RoomID)
	return err
}

const leaveRoom = `-- name: LeaveRoom :exec
delete from chat_rooms 
where user_id = $1
and room_id =  $2
`

type LeaveRoomParams struct {
	UserID uuid.UUID `json:"user_id"`
	RoomID uuid.UUID `json:"room_id"`
}

func (q *Queries) LeaveRoom(ctx context.Context, arg LeaveRoomParams) error {
	_, err := q.db.Exec(ctx, leaveRoom, arg.UserID, arg.RoomID)
	return err
}

const listPublicRooms = `-- name: ListPublicRooms :many
select room_id, name, description from rooms where room_type = 'public'
`

type ListPublicRoomsRow struct {
	RoomID      uuid.UUID   `json:"room_id"`
	Name        string      `json:"name"`
	Description pgtype.Text `json:"description"`
}

func (q *Queries) ListPublicRooms(ctx context.Context) ([]ListPublicRoomsRow, error) {
	rows, err := q.db.Query(ctx, listPublicRooms)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListPublicRoomsRow{}
	for rows.Next() {
		var i ListPublicRoomsRow
		if err := rows.Scan(&i.RoomID, &i.Name, &i.Description); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
