create table if not exists task(
    id serial constraint task_pk PRIMARY KEY,
    order_name text not null,
    start_date timestamptz
);

create table if not exists work(
    id  serial constraint work_pk PRIMARY KEY,
    task_id integer constraint work_fk REFERENCES task ON DELETE CASCADE,
    duration integer not null,
    resource integer not null,
    previous_ids integer[]
);

