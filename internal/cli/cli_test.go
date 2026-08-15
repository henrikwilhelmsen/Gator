package cli

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/config"
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/database"
)

type mockDB struct {
	users map[string]database.User
}

func (m *mockDB) CreateUser(
	ctx context.Context, arg database.CreateUserParams) (database.User, error) {
	u := database.User{
		ID:        arg.ID,
		CreatedAt: arg.CreatedAt,
		UpdatedAt: arg.UpdatedAt,
		Name:      arg.Name,
	}
	m.users[arg.Name] = u
	return u, nil
}

func (m *mockDB) GetUser(ctx context.Context, name string) (database.User, error) {
	u, ok := m.users[name]
	if !ok {
		return database.User{}, sql.ErrNoRows
	}
	return u, nil
}

func (m *mockDB) DeleteAll(ctx context.Context) error {
	m.users = map[string]database.User{}
	return nil
}

func getMockState() State {
	testUser := "jane"
	cfg := config.Config{DbURL: "postgres://example", CurrentUserName: testUser}
	db := &mockDB{
		users: make(map[string]database.User),
	}
	db.users[testUser] = database.User{Name: testUser}

	return State{
		Config: &cfg,
		Db:     db,
	}
}

// TestLogin tests that the login command sets the user to the given argument
func TestLogin(t *testing.T) {
	testUser := "john"

	// Override the config file with a tempfile
	tmpDir := t.TempDir()
	os.Setenv(config.ConfigFilePathEnvVar, tmpDir+"/testconfig.json")

	// Get a mocked state object and set up the commands
	state := getMockState()
	state.Db.CreateUser(context.Background(), database.CreateUserParams{Name: testUser})
	commands := SetupRegisterCommands()

	// Run the login command
	command := Command{Name: "login", Args: []string{testUser}}
	err := commands.Run(&state, command)
	if err != nil {
		log.Fatal(err)
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
	os.Setenv(config.ConfigFilePathEnvVar, tmpDir+"/testconfig.json")

	// Get a mocked state object and set up the commands
	state := getMockState()
	commands := SetupRegisterCommands()

	// Run the register command
	command := Command{Name: "register", Args: []string{testUser}}
	err := commands.Run(&state, command)
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
