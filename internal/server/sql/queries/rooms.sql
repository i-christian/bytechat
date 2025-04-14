-- name: InitiatePrivateRoom :one
select * from create_private_room($1, $2); 

-- name: GetPrivateRooms :many
select
    usr.first_name || ' ' || usr.last_name as full_name,
    usr.status,
    usr.user_id,
    usr.updated_at,
    shared_cr.room_id
from
    chat_rooms as specific_user_cr
join
    rooms as r on specific_user_cr.room_id = r.room_id
join
    chat_rooms as shared_cr on r.room_id = shared_cr.room_id
join
    users as usr on shared_cr.user_id = usr.user_id
where
    specific_user_cr.user_id = $1
    and r.room_type = 'private'
    and usr.user_id <> $1
order by
    usr.updated_at desc;

-- name: CreatePublicRoom :exec
insert into rooms(name, description, room_type) 
values ($1, $2, $3);

-- name: ListPublicRooms :many
select room_id, name, description from rooms where room_type = 'public';

-- name: GetRoomDetails :one
select name, room_type from rooms
where room_id = $1;

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
