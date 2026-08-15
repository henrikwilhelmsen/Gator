# Gator

A RSS feed aggre**gator** built in Go as a [boot.dev](boot.dev) guided project.

## Dev Requirements

- Go
- PostgreSQL
- Goose
- SQLC

## Dev Setup

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
