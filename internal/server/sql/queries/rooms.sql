-- name: InitiatePrivateRoom :one
select from create_private_room($1, $2); 

-- name: GetPrivateRooms :many
select 
    users.first_name || ' ' || users.last_name as full_name,
    users.status,
    users.user_id,
    users.updated_at,
    chat_rooms.room_id
from users
join chat_rooms using(user_id)
join rooms using (room_id)
where rooms.room_id = (
    select 
        rooms.room_id
    from users 
    join chat_rooms using(user_id)
    join rooms using(room_id)
    where users.user_id = $1
        and rooms.room_type = 'private'
    order by rooms.room_id
);

-- name: CreatePublicRoom :exec
insert into rooms(name, description, room_type) 
values ($1, $2, $3);

-- name: ListPublicRooms :many
select room_id, name, description from rooms where room_type = 'public';

-- name: GetRoomDetails :one
select name from rooms
where room_id = $1
and room_type = 'public';

-- name: JoinRoom :exec
insert into chat_rooms(user_id, room_id) 
values ($1, $2)
on conflict (user_id, room_id) do nothing;;

-- name: LeaveRoom :exec
delete from chat_rooms 
where user_id = $1
and room_id =  $2;

-- name: GetUsersInRoom :many
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
limit $2;
