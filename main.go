package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"

	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/cli"
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/config"
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/database"
)

func main() {
	// Load the cfg
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	// Set up the database connections
	db, err := sql.Open("postgres", cfg.DbURL)
	dbQueries := database.New(db)

	// Read arguments
	args := os.Args
	if len(args) < 2 {
		log.Fatal("Expected at least one command")
	}

	// Get the state and command objects
	state := cli.State{Config: &cfg, Db: dbQueries}
	command := cli.Command{Name: args[1], Args: args[2:]}
	commands := cli.SetupRegisterCommands()

	// Run the command
	err = commands.Run(&state, command)
	if err != nil {
		log.Fatal(err)
	}
}
