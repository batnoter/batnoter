create table if not exists default_repos
(
    id          serial primary key,
    created_at  timestamp without time zone default (now() at time zone 'utc'),
    updated_at  timestamp without time zone default (now() at time zone 'utc'),
    deleted_at  timestamp without time zone default null,
    user_id     integer not null unique,

    name                varchar(50) not null,
    visibility          varchar(20) not null,
    default_branch      varchar(50) not null,
    constraint fk_user foreign key(user_id) references users(id)
);