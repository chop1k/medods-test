create type "app"."tag_type" as enum (
    'predefined',
    'user-defined'
);

create table "app"."tags" (
    "id" integer generated always as identity,

    "name" varchar(32) not null,
    "description" text,
    "type" "app"."tag_type" not null,

    "deleted_at" timestamp,

    constraint "app_tags_id_pkey"
        primary key (id),
    constraint "app_tags_empty_chck" check (
        name != '' and (description is null or description != '')
    )
);

create table "app"."templates_tags" (
    "id" integer generated always as identity,

    "template_id" integer not null,
    "tag_id" integer not null,

    constraint "app_templates_tags_id_pkey"
        primary key (id),
    constraint "app_templates_tags_template_id_templates_id_fkey"
        foreign key (template_id)
            references "app"."templates" (id)
                on delete cascade,
    constraint "app_templates_tags_tag_id_tags_id_fkey"
        foreign key (tag_id)
            references "app"."tags" (id)
                on delete cascade
);