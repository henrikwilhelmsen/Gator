package main

import (
	"log"
	"os"

	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/cli"
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/config"
)

func main() {
	// Load the config
	config, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	// Read arguments
	args := os.Args
	if len(args) < 2 {
		log.Fatal("Expected at least one command")
	}

	// Get the state and command objects
	state := cli.State{Config: &config}
	command := cli.Command{Name: args[1], Args: args[2:]}
	commands := cli.SetupRegisterCommands()

	// Run the command
	err = commands.Run(&state, command)
	if err != nil {
		log.Fatal(err)
	}
}
