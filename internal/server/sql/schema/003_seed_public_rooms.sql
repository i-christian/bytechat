-- +goose Up
insert into rooms(name, description, room_type)
values 
    ('General', 'A public for any discussions', 'public'),
    ('Programming', 'A public room for programmers', 'public'),
    ('Music', 'A public room for discussing music', 'public'),
    ('Movie & Series', 'A public room for movies & series discussions', 'public'),
    ('Anime', 'A public room for anime discussions', 'public'),
    ('Books', 'A public room for book discussions', 'public'),
    ('Religion', 'A public room for religious discussions', 'public'),
    ('Politics', 'A public room for political discussions', 'public'),
    ('Video Games', 'A public room for video games discussions', 'public'),
    ('Sports', 'A public room for sports discussions', 'public')    
on conflict (name) do nothing;

-- +goose Down
delete from rooms;
