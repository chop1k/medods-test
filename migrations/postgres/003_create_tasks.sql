
create type "app"."task_status" as enum (
    'pending',
    'running',
    'finished',
    'cancelled',
    'moved',
    'overdue'
);

create table "app"."tasks" (
    "id" integer generated always as identity,
    "template_id" integer,
    "moved_task_id" integer,

    "status" "app"."task_status" not null,
    "notes" text,
    "date" date not null,
    "started_at" time,
    "ended_at" time,

    "deleted_at" timestamp,

    constraint "app_tasks_id_pkey"
        primary key (id),
    constraint "app_tasks_template_id_templates_id_fkey"
        foreign key (template_id)
            references "app"."templates" (id)
                on delete cascade,
    constraint "app_tasks_moved_to_task_id_fkey"
        foreign key (moved_task_id)
            references "app"."tasks" (id)
                on delete cascade
                    deferrable initially deferred,

    constraint "app_tasks_timestamps_chck" check (
        not ("started_at" is null and "ended_at" is not null)
    ),
    constraint "app_tasks_empty_chck" check (
        "notes" is null or "notes" != ''
    ),

    constraint "app_tasks_status_chck" check (
        ("status" = 'pending' and "started_at" is null and "ended_at" is null) or
        ("status" = 'running' and "started_at" is not null) or 
        ("status" = 'finished' and "started_at" is not null and "ended_at" is not null) or
        ("status" = 'cancelled' and (("started_at" is not null and "ended_at" is not null) or ("started_at" != null and "ended_at" is null))) or
        ("status" = 'moved' and "moved_task_id" is not null) or
        ("status" = 'overdue' and "started_at" is not null and "ended_at" is null)
    )
);