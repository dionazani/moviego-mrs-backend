-- 1. Person Table (The core identity)
-- sign_up, app_person, app_person, app_user
CREATE TABLE app_person (
    id UUID PRIMARY KEY,
    sign_up_from char(3), -- web (WEB), mobile (MBL),
    sign_up_status char(3), -- not yet activate (NYA), activate (ACT),
    sign_up_at TIMESTAMP,
    fullname VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE,
    mobile_phone VARCHAR(25) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

create index idx_app_person_email on app_person (email);
create index idx_app_person_mobile_phone on app_person (mobile_phone);

-- 2. User Table (Login credentials and security)
CREATE TABLE app_user (
    id UUID PRIMARY KEY,
    app_person_id UUID NOT NULL UNIQUE, -- Link to app_person
    app_user_role CHAR(3) NOT NULL,
    app_password VARCHAR(300) NOT NULL, -- Renamed for clarity
    must_change_password int default 0,
    next_change_password_date DATE DEFAULT CURRENT_DATE + INTERVAL '30 days',
    is_locked int DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (app_person_id) REFERENCES app_person (id)
);

create table app_user_activate (
	id UUID PRIMARY KEY,
	app_person_id UUID NOT NULL UNIQUE, -- Link to app_person,
	activate_by CHAR(3), -- web (WEB), mobile (MBL)
	activate_at TIMESTAMP,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (app_person_id) REFERENCES app_person (id)
);

-- drop table app_user_token;
create table app_user_token (
	id UUID primary key,
	app_user_id UUID not null,
	token_type varchar(25) default 'refresh',
	token_user varchar(200) not null,
	expire_at TIMESTAMP not null,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	foreign key (app_user_id) references app_user (id),
	unique (app_user_id, token_type)
);