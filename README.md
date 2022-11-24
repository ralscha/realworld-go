# ![RealWorld Example App with Go, Chi and SQLBoiler](logo.png)

> ### realworld-go codebase containing a real-world example (CRUD, auth, advanced patterns, etc.) that adheres to the [RealWorld](https://github.com/gothinkster/realworld) spec and API.


### [Demo](https://demo.realworld.io/)&nbsp;&nbsp;&nbsp;&nbsp;[RealWorld](https://github.com/gothinkster/realworld)


This codebase was created to demonstrate a fully-fledged backend application built with **Go, Chi, and SQLBoiler**, including CRUD operations, authentication, routing, pagination, and more.

We've gone to great lengths to adhere to the **Go** community style guides & best practices.

For more information on how this works with other frontends/backends, head over to the [RealWorld](https://github.com/gothinkster/realworld) repo.


# Develop

### Task: build tool

This project uses [Task](https://taskfile.dev/) as task runner and build tool. Checkout the installation instruction here:    
https://taskfile.dev/installation/  


### Goose: database migration

[goose](https://pressly.github.io/goose/) is used as database migration tool. The migration SQL scripts are in the [`migrations`](https://github.com/ralscha/realworld-go/tree/main/migrations) folder. The source code of the migration tool is located here: [`cmd/migrate/main.go`](https://github.com/ralscha/realworld-go/blob/main/cmd/migrate/main.go)

To build the migration tool, run `task build-goose`. The migration reads the database connection configuration from the [`app.env`](https://github.com/ralscha/realworld-go/blob/main/app.env) file.

Migration related tasks:
| Task   | Description  |
|---|---|
| `db-migration-new-go`  | Creates an empty Go migration file |
| `db-migration-new`  | Creates an empty SQL migration file |
| `db-migration-reset`  | Reverts all migrations. Results in an empty database |
| `db-migration-status`  | Shows the current applied migration version |
| `db-migration-up`  | Runs all migrations if not already applied |


### SQLBoiler: database access

The project utilizes [SQLBoiler](https://github.com/volatiletech/sqlboiler) as the database access layer. SQLBoiler depends on generated code that needs to be re-generated each time the database schema changes. SQLBoiler reads the information from the database and creates files in the folder [`internal/models`](https://github.com/ralscha/realworld-go/tree/main/internal/models). Database connection parameters are configured in the file [`sqlboiler/sqlboiler.toml`](https://github.com/ralscha/realworld-go/blob/main/sqlboiler/sqlboiler.toml). First apply the schema changes with goose and then run `task db-codegen` to generate the database code.


# Getting started

1. Install [Task](https://taskfile.dev/)
2. Start PostgreSQL with `docker compose up -d`
3. Apply database migration `task db-migration-up` 
4. Start API with `task run`
5. Run the realworld [Postman collection](https://realworld-docs.netlify.app/docs/specs/backend-specs/postman)
```
 cd postman   
 npm install   
 npm start
```