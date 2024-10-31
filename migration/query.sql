create database db_user;

create user sso_user with password 'sso123';

grant all privileges on database db_user to sso_user;

create schema if not exists sso;

create table if not exists sso.users (
	id serial primary key,
	username varchar(25) unique not null,
	name varchar(120) not null,
	dob varchar(10) not null,
	gender varchar(1) not null,
	create_date timestamp not null,
	delete_date timestamp
);

create table if not exists sso.passwords (
	user_id int references sso.users(id),
	password varchar(255) not null,
	create_date timestamp not null,
	update_date timestamp
);