package main

import (
	"log"
	"os"

	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/cli"
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/config"
)

func main() {
	config, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	// Setup and register commands
	state := cli.State{Config: &config}
	commands := cli.Commands{
		CommandsToHandlers: make(map[string]func(*cli.State, cli.Command) error),
	}
	commands.Register("login", cli.HandlerLogin)

	// Read arguments
	args := os.Args
	if len(args) < 2 {
		log.Fatal("Expected at least one command")
	}

	// Run commands
	command := cli.Command{Name: args[1], Args: args[2:]}
	err = commands.Run(&state, command)
	if err != nil {
		log.Fatal(err)
	}
}
