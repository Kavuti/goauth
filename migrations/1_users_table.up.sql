create table users (
    email varchar(255) primary key not null,
    password varchar(1000) not null,
    first_name varchar(50) not null,
    last_name varchar(50) not null,
    verified boolean not null default false
);
