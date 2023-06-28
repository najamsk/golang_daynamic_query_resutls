drop table warriors;

-- table definition
create table warriors (
    "id" serial primary key,
    "category" varchar(255) not null,
    "first_name" varchar(255) not null,
    "last_name" varchar(255),
    "teacher" varchar(255),
    "is_active" bool default false,
    "create_on" timestamptz not null default now(),
    "updated_on" timestamptz
);

insert into "public"."warriors" ("id", "category", "first_name", "last_name", "teacher", "is_active", "create_on", "updated_on") values
(1, 'ninja', 'naruto', null, 'kakashi', 't', now(), null),
(2, 'ninja', 'sasuke', null, 'kakashi', 't', now(), null),
(3, 'ninja', 'kakashi', 'hatake', 'orochimaru', 't', now(), now());

select * from warriors;
