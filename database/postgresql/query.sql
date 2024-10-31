-- name: GetUser :one
select 
    id, 
    username, 
    name, 
    email, 
    dob, 
    gender
from 
    users
where 
    username = $1;

-- name: CreateUserReturnId :one
insert into users (
    username,
    name,
    email,
    dob,
    gender,
    create_date
)
values (
    $1,
    $2,
    $3,
    $4,
    $5,
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