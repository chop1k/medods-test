create schema "app";

create table "app"."migrations" (
    "id" integer not null,

    "name" varchar(128) not null,

    "created_at" timestamp not null,

    constraint "app_migrations_id_pkey"
        primary key (id)
);