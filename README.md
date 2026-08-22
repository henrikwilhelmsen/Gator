# Gator

A RSS feed aggre**gator** built in Go as a [boot.dev](boot.dev) guided project.

## Dev Requirements

- Go
- PostgreSQL
- Goose
- SQLC

## Nix Setup (Highly Recommended)

If you use [Nix](https://nixos.org/), a fully self-contained development environment is provided via the included `flake.nix`. It automatically handles installing Go, PostgreSQL 18, Goose, SQLC, and the `boot.dev` CLI locally inside the project without modifying your global system.

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

To run the tests (does not require a database setup):

```shell
go test ./...
```

### 2. Additional Commands

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

Running migrations:

```shell
goose up
```

```shell
goose down
```

Generate database code:

```shell
sqlc generate
```

## Manual Dev Setup

### PostgreSQL

**MacOS:**

```shell
brew install postgresql@18
```

Verify the install:

```shell
psql --version
```

Start the postgres service:

```shell
brew services start postgresql@18
```

### Goose

See the [Goose install instructions](https://github.com/pressly/goose#install) or run the below command:

```shell
go install github.com/pressly/goose/v3/cmd/goose@latest
```

After installation, create a `.env` file with the appropriate data filled in, matching the [example_dot_env](./example_dot_env) file, so Goose picks up the migrations directory and you don't have to have the DB url and driver in the command.

Once the `.env` file is set up you can run the migrations with:

```shell
goose up
```

```shell
goose down
```

### SQLC

SQLC is used to generate Go code to interact with our database.

To install, see the [SQLC install docs](https://docs.sqlc.dev/en/latest/overview/install.html) or run the below command to install with Go:

```shell
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

When installed, you can generate Go code from the `./sql` directory with the `generate` command:

```shell
sqlc generate
```

See the [database package](./internal/database/db.go) for the generated code.
