-- name: InitiatePrivateRoom :one
select from create_private_room($1, $2); 

-- name: GetPrivateRooms :many
select users.first_name || ' ' || users.last_name as my_name,
    rooms.room_id
from users 
join chat_rooms on users.user_id = chat_rooms.user_id
join rooms on chat_rooms.room_id = rooms.room_id
where users.user_id = $1
and rooms.room_type = 'private'
order by rooms.room_id;

-- name: CreatePublicRoom :exec
insert into rooms(name, description, room_type) 
values ($1, $2, $3);

-- name: ListPublicRooms :many
select room_id, name, description from rooms where room_type = 'public';

-- name: JoinRoom :exec
insert into chat_rooms(user_id, room_id) 
values ($1, $2);

-- name: LeaveRoom :exec
delete from chat_rooms 
where user_id = $1
and room_id =  $2;

-- name: GetUsersInRoom :many
select 
    users.last_name || ' ' || users.first_name as full_name,
    users.status
from users
join chat_rooms using(user_id)
join rooms using (room_id)
where rooms.room_id = $1;
