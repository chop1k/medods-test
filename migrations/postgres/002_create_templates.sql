create table "app"."templates" (
    "id" integer generated always as identity,

    "name" varchar(64) not null,
    "description" text,
    "starts_at" time not null,
    "ends_at" time not null,

    "enabled" boolean not null,

    "scheduling" jsonb not null,

    constraint "app_templates_id_pkey"
        primary key (id),
    constraint "app_templates_empty_chck" check (
        "name" != '' and ("description" is null or "description" != '')
    )
);