create table mst_user_role (
	id UUID,
	code char(5) unique not null,
	names varchar(25) not null,
	description varchar(50) not null,
	primary key (id)
);

-- 1. Person Table (The core identity)
-- sign_up, app_person, app_person, app_user
CREATE TABLE app_person (
    id UUID PRIMARY KEY,
    sign_up_from char(3), -- web (WEB), mobile (MBL),
    sign_up_at TIMESTAMP,
    fullname VARCHAR(100) NOT NULL,
    gender CHAR(1) NOT NULL,
    email VARCHAR(255) unique not null,
    mobile_phone VARCHAR(25) unique not null,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

create index idx_app_person_email on app_person (email);
create index idx_app_person_mobile_phone on app_person (mobile_phone);

-- 2. User Table (Login credentials and security)
CREATE TABLE app_user (
    id UUID PRIMARY KEY,
    app_person_id UUID NOT NULL UNIQUE, -- Link to app_person
    mst_user_role_id UUID NOT NULL,
    app_password VARCHAR(300) NOT NULL, -- Renamed for clarity
    must_change_password int default 0,
    next_change_password_date DATE DEFAULT CURRENT_DATE + INTERVAL '30 days',
    is_locked int DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    failed_attempt_count INT DEFAULT 0,
    lockout_until DATETIME DEFAULT null,
    FOREIGN KEY (app_person_id) REFERENCES app_person (id),
    FOREIGN KEY (mst_user_role_id) REFERENCES mst_user_role (id)
);

-- alter table app_user add failed_attempt_count INT DEFAULT 0;
-- alter table app_user add lockout_until TIMESTAMP DEFAULT NULL;

-- drop table app_user_token;
create table app_user_token (
	id UUID primary key,
	app_user_id UUID not null,
	token_type varchar(25) default 'refresh',
	token_user varchar(500) not null,
	expire_at TIMESTAMP not null,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	foreign key (app_user_id) references app_user (id),
	unique (app_user_id, token_type)
);