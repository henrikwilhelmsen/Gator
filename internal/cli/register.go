package cli

import "git.hwanimation.tech/henrikwilhelmsen/gator/internal/state"

// SetupRegisterCommands sets up and returns a Commands struct, with all supported
// commands registered.
func SetupRegisterCommands() *Commands {
	commands := Commands{
		CommandsToHandlers: make(map[string]func(*state.State, Command) error),
	}
	commands.Register("login", handlerLogin)
	commands.Register("register", handlerRegister)
	commands.Register("reset", handlerReset)
	commands.Register("users", handlerListUsers)
	commands.Register("agg", handlerAgg)
	commands.Register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commands.Register("feeds", handlerListFeeds)
	commands.Register("follow", middlewareLoggedIn(handlerFollow))
	commands.Register("following", middlewareLoggedIn(handlerFollowing))
	return &commands
}
