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
