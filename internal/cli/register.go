package cli

import "git.hwanimation.tech/henrikwilhelmsen/gator/internal/state"

// SetupRegisterCommands sets up and returns a Commands struct, with all supported
// commands registered. The current commands are: login
func SetupRegisterCommands() *Commands {
	commands := Commands{
		CommandsToHandlers: make(map[string]func(*state.State, Command) error),
	}
	commands.Register("login", HandlerLogin)
	commands.Register("register", HandlerRegister)
	commands.Register("reset", HandlerReset)
	commands.Register("users", HandlerListUsers)
	commands.Register("agg", HandlerAgg)
	commands.Register("addfeed", HandlerAddFeed)
	commands.Register("feeds", HandlerListFeeds)
	return &commands
}
