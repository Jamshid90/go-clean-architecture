CREATE TABLE IF NOT EXISTS "refresh_token" (
    "user_id" character varying(20) NOT NULL,
    "token" character varying(500) NOT NULL,
    "created_at" timestamp(0) without time zone);