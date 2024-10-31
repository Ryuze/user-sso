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
desc;