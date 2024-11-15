-- name: CreateUserReturnId :one
insert into users (
    username,
    name,
    email,
    dob,
    gender,
    allowed_services,
    create_date
)
values (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    now()
)
returning id;

-- name: CreatePasswordReturnUserId :one
insert into passwords (
    user_id,
    password,
    create_date
)
values (
    $1,
    $2,
    now()
)
returning user_id;

-- name: GetUserLatestPassword :one
select 
    users.id,
    users.username,
    passwords.password
from 
    users
join 
    passwords
on 
    users.id = passwords.user_id
where 
    users.username = $1
and 
    passwords.update_date is null
order by
    passwords.create_date
desc
limit 1;

-- name: GetUser :one
select 
    id,
    username,
    name,
    allowed_services
from 
    users
where 
    username = $1;

-- name: GetService :one
select 
    service_name,
    public_key,
    status
from 
    services
where 
    service_name = $1;