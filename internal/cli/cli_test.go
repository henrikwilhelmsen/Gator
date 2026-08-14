package cli

import (
	"log"
	"os"
	"testing"

	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/config"
)

// TestLogin tests that the login command sets the user to the given argument
func TestLogin(t *testing.T) {
	testUser := "john"
	tmpDir := t.TempDir()
	os.Setenv(config.ConfigFilePathEnvVar, tmpDir+"/testconfig.json")

	cfg := config.Config{DbURL: "postgres://example", CurrentUserName: "jane"}
	state := State{Config: &cfg}
	commands := SetupRegisterCommands()

	command := Command{Name: "login", Args: []string{testUser}}
	err := commands.Run(&state, command)
	if err != nil {
		log.Fatal(err)
	}

	got, err := config.Read()
	if got.CurrentUserName != testUser || err != nil {
		t.Fatalf("login john = %v, %v, want match for %q, nil", got.CurrentUserName, err, testUser)
	}
}
