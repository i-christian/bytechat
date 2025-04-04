-- +goose Up

-- +goose StatementBegin
create or replace function create_private_room(user_a uuid, user_b uuid)
    returns uuid
    language plpgsql
as
$$
declare
    room_id_var uuid;
    type_of_room text;
    room_name_var text;
    ordered_user_id_1 uuid;
    ordered_user_id_2 uuid;
begin
    type_of_room := 'private';

    -- Ensure consistent ordering of user IDs for the room name
    if user_a < user_b then
        ordered_user_id_1 := user_a;
        ordered_user_id_2 := user_b;
    else
        ordered_user_id_1 := user_b;
        ordered_user_id_2 := user_a;
    end if;

    room_name_var := ordered_user_id_1::text || '_' || ordered_user_id_2::text;

    insert into rooms(name, room_type)
        values (room_name_var, type_of_room)
        on conflict(name) do nothing
        returning room_id into room_id_var;

    if room_id_var is null then
        select room_id into room_id_var
        from rooms
        where name = room_name_var;
    end if;

    insert into chat_rooms (user_id, room_id)
    values 
        (user_a, room_id_var),
        (user_b, room_id_var)
    on conflict (user_id, room_id) do nothing;

    return room_id_var;
end;
$$;
-- +goose StatementEnd

-- +goose StatementBegin
create or replace function fn_update_timestamp()
    returns trigger
    language plpgsql
as 
$$
begin
  new.updated_at = CURRENT_TIMESTAMP;
  return new;
end;
$$;
-- +goose StatementEnd

create trigger trg_update_users_timestamp
before update on users
for each row
execute function fn_update_timestamp();

-- +goose Down
drop function if exists func_name();
drop trigger if exists trg_update_users_timestamp on users;
drop function if exists fn_update_timestamp();

