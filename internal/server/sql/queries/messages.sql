-- name: CreateMessage :one
insert into messages(user_id, room_id, text) 
values ($1, $2, $3)
returning *;

-- name: ListMessagesByRoom :many
select m.message_id, m.room_id, m.user_id, m.text, m.created_at, u.first_name, u.last_name
from messages m
join users u on m.user_id = u.user_id
where m.room_id = $1
order by m.created_at ASC
limit $2;
