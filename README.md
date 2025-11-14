# health-check

#### Run the app and tests
```
>>> go run cmd/main.go --config sample-config.json
>>> go test ./... -v
```

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

#### Create table and grant permissions
```
CREATE TABLE check_result_db (
    id SERIAL PRIMARY KEY,
    api_id INTEGER NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    status_code INTEGER NOT NULL,
    success BOOLEAN NOT NULL,
    response_time_ms BIGINT NOT NULL,
    response_snippet TEXT,
    error_message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

>>> grant select, insert, update, delete on table check_result_db to admin;
>>> grant usage, select on sequence check_result_db_id_seq to admin;

```

#### register API 
```
 POST  http://localhost:8080/api/v1/register/

 {
  "name": "User Service Health",
  "url": "https://jsonplaceholder.typicode.com/posts/1",
  "method": "GET",
  "headers": {
    "Accept": "application/json"
  },
  "body": "hello again",
  "interval_seconds": 3,
  "enabled": true,
  "webhook": {
    "url": "https://webhook.site/your-custom-url",
    "headers": {
      "Authorization": "Bearer my-secret-token"
    }
  }
}


{
  "name": "Failing API",
  "url": "http://localhost:9999/nonexistent", 
  "method": "GET",
  "interval_seconds": 10,
  "enabled": true,
  "webhook": {
  "url": "http://localhost:9000/webhook",
  "headers": {
    "Content-Type": "application/json"
  }
}
}
```

#### health-check API
```
POST  http://localhost:8080/api/v1/:api_id/start

```

#### list api
```
GET  http://localhost:8080/api/v1/apis
```

#### delete api
```
DELETE  http://localhost:8080/api/v1/delete/:api_id
```