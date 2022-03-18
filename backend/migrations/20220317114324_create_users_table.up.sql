create table if not exists users
(
    id              serial primary key,
    created_at      timestamp without time zone default (now() at time zone 'utc'),
    updated_at      timestamp without time zone default (now() at time zone 'utc'),
    deleted_at      timestamp without time zone default null,
    disabled_at     timestamp without time zone default null,
    email           varchar(255) not null unique,
    name            varchar(255) null,
    location        varchar(255) null,
    avatar_url      varchar(255) null,
    github_id       int null,
    github_username varchar(255) null,
    github_token    varchar(255) not null
);