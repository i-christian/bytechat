-- name: CreateRoom :exec
insert into rooms(name, description, room_type) 
values ($1, $2, $3);
