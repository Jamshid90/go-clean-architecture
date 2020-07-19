CREATE TABLE IF NOT EXISTS "user" (
    "id" serial primary key,
    "email" character varying(50) NOT NULL,
    "status" character varying(50) NOT NULL,
    "first_name" character varying(50) NOT NULL,
    "last_name" character varying(50) NOT NULL,
    "password" character varying(200) NOT NULL,
    "created_at" timestamp(0) without time zone,
    "updated_at" timestamp(0) without time zone);