{
  description = "Gator RSS aggregator development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      supportedSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forEachSupportedSystem = f: nixpkgs.lib.genAttrs supportedSystems (system: f {
        pkgs = nixpkgs.legacyPackages.${system};
      });
    in
    {
      devShells = forEachSupportedSystem ({ pkgs }:
        let
          pg-start = pkgs.writeShellApplication {
            name = "pg-start";
            runtimeInputs = [ pkgs.postgresql_18 ];
            text = ''
              if pg_ctl status >/dev/null 2>&1; then
                echo "PostgreSQL is already running."
              else
                echo "Starting PostgreSQL..."
                pg_ctl -l "$PGDATA/server.log" start
                until pg_isready -h "$PGHOST" -p "$PGPORT" >/dev/null 2>&1; do
                  sleep 0.1
                done
                if ! psql -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" -lqt | cut -d \| -f 1 | grep -qw "$PGDATABASE"; then
                  echo "Creating database '$PGDATABASE'..."
                  createdb -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" "$PGDATABASE"
                fi
              fi
            '';
          };

          pg-stop = pkgs.writeShellApplication {
            name = "pg-stop";
            runtimeInputs = [ pkgs.postgresql_18 ];
            text = ''
              if pg_ctl status >/dev/null 2>&1; then
                echo "Stopping PostgreSQL..."
                pg_ctl stop
              else
                echo "PostgreSQL is not running."
              fi
            '';
          };

          pg-console = pkgs.writeShellApplication {
            name = "pg-console";
            runtimeInputs = [ pkgs.postgresql_18 ];
            text = ''
              psql -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" -d "$PGDATABASE"
            '';
          };
        in
        {
          default = pkgs.mkShell {
            packages = with pkgs; [
              go
              golangci-lint
              gitlab-ci-ls
              nil # nix lsp
              goose
              sqlc
              postgresql_18
              fish
              pg-start
              pg-stop
              pg-console
            ];

            shellHook = ''
            # Contain the database files in a .nix directory
            export PGDATA="$PWD/.nix/db"
            export PGPORT=5432
            export PGHOST="localhost"
            export PGDATABASE="gator"
            export PGUSER="$(whoami)"

            # Isolate Go binaries so bootdev is local to the project
            export GOBIN="$PWD/.nix/bin"
            export PATH="$GOBIN:$PATH"

            # Isolate Gator's config file within the project .nix directory
            export GATOR_CONFIG_FILEPATH="$PWD/.nix/.gatorconfig.json"

            # 1. Initialize local PostgreSQL database if it doesn't exist
            if [ ! -d "$PGDATA" ]; then
              echo "Initializing local PostgreSQL database in $PGDATA..."
              initdb --no-locale -U "$PGUSER"

              # Adjust pg_hba.conf to allow trust (passwordless) connections locally
              echo "host all all 127.0.0.1/32 trust" >> "$PGDATA/pg_hba.conf"
              echo "host all all ::1/128 trust" >> "$PGDATA/pg_hba.conf"

              # Configure Postgres to listen only on localhost
              echo "listen_addresses = 'localhost'" >> "$PGDATA/postgresql.conf"
              echo "port = $PGPORT" >> "$PGDATA/postgresql.conf"
            fi

            # 2. Automatically generate .gatorconfig.json if missing
            if [ ! -f "$GATOR_CONFIG_FILEPATH" ]; then
              echo "Initializing $GATOR_CONFIG_FILEPATH..."
              cat <<EOF > "$GATOR_CONFIG_FILEPATH"
{
  "db_url": "postgres://$PGUSER@localhost:$PGPORT/$PGDATABASE?sslmode=disable",
  "current_user_name": ""
}
EOF
            fi

            # 3. Automatically generate Goose .env configuration if missing
            if [ ! -f ".env" ]; then
              echo "Initializing .env for Goose..."
              cat <<EOF > .env
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://$PGUSER@localhost:$PGPORT/$PGDATABASE?sslmode=disable
GOOSE_MIGRATION_DIR=./sql/schema
EOF
            fi

            # 4. Automatically install boot.dev CLI (bootdev) if missing
            if ! command -v bootdev &> /dev/null; then
              echo "Installing boot.dev CLI (bootdev) locally to $GOBIN..."
              mkdir -p "$GOBIN"
              go install github.com/bootdotdev/bootdev@latest
            fi

            # 5. pg-start, pg-stop, and pg-console are provided as standalone
            # packages (see the `let` block above) so they work under any shell.

            echo "=== Gator Nix Dev Environment ==="
            echo "Available tools:"
            echo "  - go: $(go version)"
            echo "  - goose: $(goose --version 2>&1 | head -n 1)"
            echo "  - sqlc: $(sqlc version)"
            echo "  - postgresql: $(postgres --version)"
            echo "  - bootdev: $(bootdev --version 2>/dev/null || echo "Installed at $GOBIN/bootdev")"
            echo ""
            echo "Helper Commands:"
            echo "  pg-start    Start the local PostgreSQL server and create '$PGDATABASE' DB"
            echo "  pg-stop     Stop the local PostgreSQL server"
            echo "  pg-console  Connect directly to the database using psql"
            echo "============================="

            if [ -z "$IN_GATOR_DEVSHELL" ]; then
              export IN_GATOR_DEVSHELL=1
              exec fish
            fi

            '';
          };
        });
    };
}
