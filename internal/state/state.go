package state

import (
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/config"
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/database"
)

type State struct {
	Config *config.Config
	Db     database.Querier
}
