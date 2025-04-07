-- name: CreateMessage :one
insert into messages(user_id, room_id, text) 
values ($1, $2, $3)
returning *;
