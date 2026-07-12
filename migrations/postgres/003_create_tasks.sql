
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
    "started_at" timestamp,
    "ended_at" timestamp,

    constraint "app_tasks_id_pkey"
        primary key (id),
    constraint "app_tasks_template_id_templates_id_fkey"
        foreign key (template_id)
            references "app"."templates" (id)
                on delete cascade,
    constraint "app_tasks_moved_to_task_id_fkey"
        foreign key (moved_task_id)
            references "app"."tasks" (id)
                on delete cascade,

    constraint "app_tasks_timestamps_chck" check (
        not ("started_at" is null and "ended_at" is not null)
    ),
    constraint "app_tasks_empty_chck" check (
        "notes" is null or "notes" != ''
    ),
    constraint "app_tasks_one_of_fkeys_chck" check (
        "status" != 'moved' and "template_id" is not null
    ),

    constraint "app_tasks_pending_chck" check (
        "status" = 'pending' and "started_at" is null and "ended_at" is null
    ),
    constraint "app_tasks_running_chck" check (
        "status" = 'running' and "started_at" is not null
    ),
    constraint "app_tasks_finished_chck" check (
        "status" = 'finished' and "started_at" is not null and "ended_at" is not null
    ),
    constraint "app_tasks_cancelled_chck" check (
        "status" = 'cancelled' and (
            ("started_at" is not null and "ended_at" is not null) or ("started_at" != null and "ended_at" is null)
        )
    ),
    constraint "app_tasks_moved_chck" check (
        "status" = 'moved' and "moved_task_id" is not null and "template_id" is null
    ),
    constraint "app_tasks_overdue_chck" check (
        "status" = 'overdue' and "started_at" is not null and "ended_at" is null
    )
);