package cli

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/config"
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/database"
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/mock"
	"github.com/google/uuid"
)

// TestLogin tests that the login command sets the user to the given argument
func TestLogin(t *testing.T) {
	testUser := "john"

	// Override the config file with a tempfile
	tmpDir := t.TempDir()
	err := os.Setenv(config.ConfigFilePathEnvVar, tmpDir+"/testconfig.json")
	if err != nil {
		t.Fatalf("unexpected error when setting config env var: %v", err)
	}

	// Get a mocked state object and set up the commands
	state := mock.GetMockState()
	_, err = state.Db.CreateUser(context.Background(), database.CreateUserParams{Name: testUser})
	if err != nil {
		t.Fatalf("unexpected error when creating test user: %v", err)
	}
	commands := SetupRegisterCommands()

	// Run the login command
	command := Command{Name: "login", Args: []string{testUser}}
	err = commands.Run(&state, command)
	if err != nil {
		t.Fatalf("unexpected error running login command: %v", err)
	}

	// Read the config and check that CurrentUserName has been updated.
	got, err := config.Read()
	if got.CurrentUserName != testUser || err != nil {
		t.Fatalf("login john = %v, %v, want match for %q, nil",
			got.CurrentUserName, err, testUser)
	}
}

// TestRegister tests that the register command registers a user in the database
// and sets them as the current user in the config
func TestRegister(t *testing.T) {
	testUser := "alice"

	// Override the config file with a tempfile
	tmpDir := t.TempDir()
	err := os.Setenv(config.ConfigFilePathEnvVar, tmpDir+"/testconfig.json")
	if err != nil {
		t.Fatalf("unexpected error when setting config env var: %v", err)
	}

	// Get a mocked state object and set up the commands
	state := mock.GetMockState()
	commands := SetupRegisterCommands()

	// Run the register command
	command := Command{Name: "register", Args: []string{testUser}}
	err = commands.Run(&state, command)
	if err != nil {
		t.Fatalf("unexpected error running register command: %v", err)
	}

	// Verify the user was registered in the database
	user, err := state.Db.GetUser(context.Background(), testUser)
	if err != nil {
		t.Fatalf("expected user %q to be in the database, but got: %v", testUser, err)
	}
	if user.Name != testUser {
		t.Fatalf("expected registered user name to be %q, got %q", testUser, user.Name)
	}

	// Verify current user config was updated
	got, err := config.Read()
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}
	if got.CurrentUserName != testUser {
		t.Fatalf("register alice = %q, want current user config match for %q",
			got.CurrentUserName, testUser)
	}
}

// TestReset tests that the reset command deletes all users in the database
func TestReset(t *testing.T) {
	// Override the config file with a tempfile
	tmpDir := t.TempDir()
	err := os.Setenv(config.ConfigFilePathEnvVar, tmpDir+"/testconfig.json")
	if err != nil {
		t.Fatalf("unexpected error when setting config env var: %v", err)
	}

	// Get a mocked state object and set up the commands
	state := mock.GetMockState()
	commands := SetupRegisterCommands()

	// Register a random number of users
	for i := range rand.Intn(20) {
		userName := fmt.Sprintf("testUser%d", i)
		cmdRegister := Command{Name: "register", Args: []string{userName}}
		err = commands.Run(&state, cmdRegister)
		if err != nil {
			t.Fatalf("unexpected error when creating test user: %v", err)
		}
	}

	// Run the reset cmdReset
	cmdReset := Command{Name: "reset", Args: []string{}}
	err = commands.Run(&state, cmdReset)
	if err != nil {
		t.Fatalf("failed to run the reset command: %v", err)
	}

	// Check that the database has been cleared
	got, err := state.Db.GetUsers(context.Background())
	if err != nil {
		t.Fatalf("failed to read users from db: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("database returned %d users, expected 0 after reset command", len(got))
	}
}

// TestAddFeed tests that the addfeed command adds a feed to the database with the
// expected name and url.
func TestAddFeed(t *testing.T) {
	// Override the config file with a tempfile
	tmpDir := t.TempDir()
	err := os.Setenv(config.ConfigFilePathEnvVar, tmpDir+"/testconfig.json")
	if err != nil {
		t.Fatalf("unexpected error when setting config env var: %v", err)
	}

	testUser := "jill"
	testUserID := uuid.New()

	// Get a mocked state object and set up the user data
	state := mock.GetMockState()

	_, err = state.Db.CreateUser(
		context.Background(),
		database.CreateUserParams{Name: testUser, ID: testUserID},
	)
	if err != nil {
		t.Fatalf("unexpected error when creating test user: %v", err)
	}

	err = state.Config.SetUser(testUser)
	if err != nil {
		t.Fatalf("unexpected error when setting current user: %v", err)
	}

	// Set up the commands
	commands := SetupRegisterCommands()

	// Run the addfeed command
	feedName := "HW Animation Tech Blog"
	feedUrl := "https://hwanimation.tech/feed"
	cmdReset := Command{Name: "addfeed", Args: []string{feedName, feedUrl}}
	err = commands.Run(&state, cmdReset)
	if err != nil {
		t.Fatalf("failed to run the reset command: %v", err)
	}

	// Check that the feed was added to the database
	got, err := state.Db.GetFeed(context.Background(), feedUrl)
	if err != nil {
		t.Fatalf("failed to get feed from db: %v", err)
	}
	if got.Name != feedName || got.Url != feedUrl || got.UserID != testUserID {
		t.Fatalf(`
			"added feed data mismatch,
			got name: %s url: %s user_id: %s,
			want name: %s url: %s, user_id: %s"`,
			got.Name, got.Url, got.UserID,
			feedName, feedUrl, testUserID)
	}
}
