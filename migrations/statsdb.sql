create table route
(
    name   varchar(255)         not null,
    load   numeric              not null,
    cargo  varchar(100)         not null,
    actual boolean default true not null,
    id     bigserial
        primary key
);

alter table route
    owner to postgres;