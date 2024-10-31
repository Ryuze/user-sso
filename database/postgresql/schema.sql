create table users (
	id serial primary key,
	username varchar(25) unique not null,
	name varchar(120) not null,
	email varchar(50) not null,
	dob varchar(10) not null,
	gender varchar(15) not null,
	create_date timestamp not null,
	delete_date timestamp
);

create table passwords (
	user_id int references users(id),
	password varchar(255) not null,
	create_date timestamp not null,
	update_date timestamp
);