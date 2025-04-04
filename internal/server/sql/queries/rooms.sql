-- name: CreateRoom :exec
insert into rooms(name, description, room_type) 
values ($1, $2, $3);

-- name: ListPublicRooms :many
select name, description from rooms;

-- name: JoinRoom :exec
insert into chat_rooms(user_id, room_id) 
values ($1, $2);

-- name: LeaveRoom :exec
delete from chat_rooms 
where user_id = $1
and room_id =  $2;

-- name: ListUsersInRoom :many
select 
    users.last_name || ' ' || users.first_name as full_name,
    users.status
from users
join chat_rooms using(user_id)
join rooms using (room_id)
where rooms.room_id = $1;
