create table if not exists notes
(
    id          serial primary key,
    created_at  timestamp without time zone default (now() at time zone 'utc'),
    updated_at  timestamp without time zone default (now() at time zone 'utc'),
    deleted_at  timestamp without time zone default null,
    email       varchar(255) not null,
    title       varchar(255) not null,
    content     text not null
);