begin;

create table if not exists machines (
    "name" text not null,
    "token" text not null,
    "address" text not null,
    "updated" bigint not null default date_part('epoch'::text, now()),
    primary key ("name")
);

commit;