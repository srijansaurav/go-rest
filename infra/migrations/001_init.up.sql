-- function to set `created_at` field as current time
create or replace function set_created_at_column()
returns trigger as $$
begin
    NEW.created_at = now();
    return NEW;
end;
$$ language plpgsql;

-- new table `user`
create table auth_user (
    -- auto generated fields
    id             bigserial primary key,
    created_at     timestamp with time zone not null,
    -- user defined fields
    username       varchar (16) not null,
    password       varchar (512) not null
);

-- create indexes on auth_user table
create unique index idx_auth_user_username on auth_user (username);

-- trigger to set `created_at` field on every insert in auth_user table
create trigger auth_user_set_created_at
    before insert on auth_user
    for each row execute procedure set_created_at_column();
