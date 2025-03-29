
create table public.tbl_topic
(
    topic       varchar(255) not null
        constraint tbl_topic_pk
            primary key,
    pod_run integer,
    script_version int ,
    script_name varchar(255),
    in_active    varchar(1) default 'N',
    created_at  timestamp,
    create_by   varchar(100),
    update_at   timestamp,
    update_by   varchar(100)
);

create table public.tbl_script
(
    id varchar(255) primary key ,
    script_key varchar(255),
    script_name varchar(255),
    version int default 0,
    is_deleted    varchar(1) default 'N',
    created_at  timestamp,
    create_by   varchar(100)
);
CREATE INDEX idx_tbl_script_script_name_version
    ON public.tbl_script (script_name,version);

create table public.tbl_connection_pool
(
    id serial primary key ,
    name varchar(255),
    type varchar(255),
    end_point varchar(255),
    username varchar(255),
    password varchar(255),
    key TEXT,
    is_deleted    varchar(1) default 'N',
    created_at  timestamp,
    create_by   varchar(100),
    update_at   timestamp,
    update_by   varchar(100),
    unique (name,type)
);
CREATE INDEX idx_tbl_connection_pool_name_type
    ON public.tbl_connection_pool (name,type);
CREATE INDEX idx_tbl_connection_pool_type
    ON public.tbl_connection_pool (type);


create table public.tbl_topic_connection
(
    topic       varchar(255),
    connection_id integer,
    is_deleted    varchar(1) default 'N',
    created_at  timestamp,
    create_by   varchar(100),
    update_at   timestamp,
    update_by   varchar(100),
    primary key (topic,connection_id)
);

--INSERT INTO pod_registry (topic, pod_name, running_at,created_at) VALUES ($1, $2, NOW(), NOW())`, topicName, podName
-- select running_at between now() - config_time and batch update running_at in config_time-config_time_renew
-- batch eod delete running_at < (now() - config_time + 1hour) clean data


create table public.tbl_pod_registry
(
    topic       varchar(255),
    pod_name       varchar(255),
    is_deleted    varchar(1) default 'N',
    created_at  timestamp,
    running_at timestamp,
    primary key(topic,pod_name)
);
CREATE INDEX idx_tbl_pod_registry_pod_name
    ON public.tbl_pod_registry (pod_name);
CREATE INDEX idx_tbl_pod_registry_running_at
    ON public.tbl_pod_registry (running_at);

CREATE INDEX idx_tbl_pod_registry_pod_name_running_at
    ON public.tbl_pod_registry (pod_name,running_at);



---

create table public.tbl_user
(
    username   varchar(255) primary key,
    password   varchar(255),
    role       varchar(255),
    is_deleted varchar(1) default 'N',
    created_at timestamp,
    create_by  varchar(100),
    update_at  timestamp,
    update_by  varchar(100)
);


create table public.tbl_role
(
    role_name  varchar(255) primary key,
    is_deleted varchar(1) default 'N',
    created_at timestamp,
    create_by  varchar(100),
    update_at  timestamp,
    update_by  varchar(100)
);


create table public.tbl_permission
(
    role_name  varchar(255),
    api_path   varchar(255),
    is_deleted varchar(1) default 'N',
    created_at timestamp,
    create_by  varchar(100),
    update_at  timestamp,
    update_by  varchar(100)
);
