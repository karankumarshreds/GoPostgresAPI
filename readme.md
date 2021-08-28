### Packages to get mux and pg driver:

`go get -u github.com/gorilla/mux`

`go get -u github.com/lib/pq`

https://dev.to/techschoolguru/how-to-write-run-database-migration-in-golang-5h6g

Guide:

https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql

## MIGRATIONS

Source: https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md

-> Download the `migrate` command line for windows from:
https://github.com/golang-migrate/migrate/releases/tag/v4.14.1

-> Set the `POSTGRESQL_URL` environment variable for `migrate` cli to create migrations for us.
`export POSTGRESQL_URL='postgres://mux:password@localhost:5432/mux_db?sslmode=disable'`

-> Now we need to create `up` and `down` files for a table that our application needs. This CLI tool will create it for us

`mkdir migrations && cd migrations`

`migrate create -ext sql -dir ./ -seq create_products_table`

This will generate two files:

`000001_create_products_table.up.sql` & `000001_create_products_table.up.sql`

-> Now add SQL in these files (create table and drop table sql)

-> Now use `migrate` CLI to run these migrations.

`migrate -database ${POSTGRESQL_URL} -path <path_to_migration_files> up`

OR `migrate -database ${POSTGRESQL_URL} -path <path_to_migration_files> down`

-> Now check if the table was created:

`psql mux_db -c "\d products"`

`psql <db_name> -c "\d <table_name>"`
