drop table warriors;

-- Table Definition
CREATE TABLE warriors (
    "id" serial primary KEY,
    "category" varchar(255) NOT NULL,
    "first_name" varchar(255) NOT NULL,
    "last_name" varchar(255),
    "teacher" varchar(255),
    "is_active" bool DEFAULT false,
    "create_on" timestamptz NOT NULL DEFAULT now(),
    "updated_on" timestamptz
);

INSERT INTO "public"."warriors" ("id", "category", "first_name", "last_name", "teacher", "is_active", "create_on", "updated_on") VALUES
(1, 'ninja', 'naruto', NULL, 'kakashi', 't', now(), NULL),
(2, 'ninja', 'sasuke', NULL, 'kakashi', 't', now(), NULL),
(3, 'ninja', 'kakashi', 'hatake', 'orochimaru', 't', now(), now());

select * from warriors;
