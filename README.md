# health-check

#### Create a database connection
``` 
> sudo -i -u postgres
> psql
> create database "health-check";
> CREATE USER admin WITH PASSWORD '12qw!@QW';
> grant all privileges on database "health-check" to admin;
> grant usage, create on schema public to admin;
> alter default privileges in schema public grant all on tables to admin;
```

#### Connect to database
> psql -h localhost -p 5432 -U admin -d "health-check";

#### Create table and grant permissions
```
CREATE TABLE monitored_apis (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    url TEXT NOT NULL,
    method TEXT NOT NULL,
    headers JSONB,
    body TEXT,
    interval_seconds BIGINT NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    last_status TEXT,
    last_checked_at TIMESTAMPTZ,
    webhook_url TEXT NOT NULL,
    webhook_headers JSONB
);

>>> sudo -i -u postgres
>>> psql -U postgres -d "health-check"
>>> grant select, insert, update, delete on table monitored_apis to admin;
>>> grant usage, select on sequence monitored_apis_id_seq to admin;
```
