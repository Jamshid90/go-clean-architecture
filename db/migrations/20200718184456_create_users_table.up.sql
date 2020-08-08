CREATE TABLE IF NOT EXISTS "user" (
    "id" character varying(20) NOT NULL,
    "email" character varying(50) DEFAULT '',
    "phone" character varying(20) DEFAULT '',
    "gender" character varying(20) DEFAULT '',
    "status" character varying(50) NOT NULL,
    "first_name" character varying(50) NOT NULL,
    "last_name" character varying(50) NOT NULL,
    "password" character varying(200) NOT NULL,
    "birth_date" date DEFAULT CURRENT_DATE,
    "created_at" timestamp(0) without time zone DEFAULT CURRENT_DATE,
    "updated_at" timestamp(0) without time zone DEFAULT CURRENT_DATE,
    CONSTRAINT user_pkey PRIMARY KEY (id));







