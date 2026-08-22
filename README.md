# Gator

A RSS feed aggre**gator** built in Go as a [boot.dev](boot.dev) guided project.

## Dev Setup

This project uses [Nix](https://nixos.org/) to manage the dev environment, install it on your machine and then follow the quick-start.

The following requirements will be installed and configured automatically:

- [Go](https://go.dev/)
- [PostgreSQL](https://www.postgresql.org/)
- [Goose](https://github.com/pressly/goose)
- [SQLC](https://sqlc.dev)

See the [nix flake](./flake.nix) for the full configuration.

### 1. Quick-Start

Run the following command in the project root to activate the environment:

```shell
nix develop
```

Start the database and run migrations:

```shell
pg-start
goose up
```

You can now run the CLI:

```shell
go run . register <username>
```

To run the tests:

```shell
go test ./...
```

### 2. Additional Commands

#### PostgreSQL

Start the local, database server:

```shell
pg-start
```

To connect to the database command line:

```shell
pg-console
```

To stop the PostgreSQL server when you are done:

```shell
pg-stop
```

#### Goose

Goose is used to run migrations, with the following commands:

```shell
goose up
```

```shell
goose down
```

#### SQLC

SQLC is used to generate Go code to interact with our database, from the queries and migrations in the [sql directory](./sql).

To generate database code:

```shell
sqlc generate
```

The generated code ends up in [./internal/database](./internal/database)
