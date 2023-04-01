create table workers
(
    id            serial primary key,
    name          varchar(50)  not null,
    surname       varchar(50)  not null,
    fathers_name  varchar(50)  not null,
    phone         varchar(25)  not null unique,
    role          varchar(50)  not null default 'worker',
    password_hash varchar(255) not null
);

create table documents
(
    id          serial primary key,
    car         varchar(50)                  not null,
    car_id      varchar(50)                  not null unique,
    waybill     int                          not null,
    driver_name varchar(255)                 not null,
    gas_amount  int check ( gas_amount > 0 ) not null,
    gas_type    varchar(50)                  not null,
    issue_date  date                         not null
);

create table workers_documents
(
    id          serial primary key,
    worker_id   int references workers (id) on delete cascade   not null,
    document_id int references documents (id) on delete cascade not null
);