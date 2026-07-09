
create type "app"."task_status" as enum (
    'pending',
    'running',
    'finished',
    'cancelled'
);

create table "app"."tasks" (
    "id" integer generated always as identity,
    "template_id" integer not null,

    "status" "app"."task_status" not null,
    "notes" text,
    "started_at" timestamp,
    "ended_at" timestamp,

    constraint "app_tasks_id_pkey"
        primary key (id),
    constraint "app_tasks_template_id_templates_id_fkey"
        foreign key (template_id)
            references "app"."templates" (id)
                on delete cascade,
    constraint "app_tasks_timestamps_chck" check (
        not (started_at is null and ended_at is not null)
    ),
    constraint "app_tasks_empty_chck" check (
        notes is null or notes != ''
    )
);