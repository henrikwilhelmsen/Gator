// Copyright 2026 Henrik Wilhelmsen. All rights reserved.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/config"
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/database"
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/state"
	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"github.com/stapelberg/postgrestest"

	_ "github.com/lib/pq"
)

var pgt *postgrestest.Server

func TestMain(m *testing.M) {
	var err error
	pgt, err = postgrestest.Start(context.Background())
	if err != nil {
		panic(err)
	}
	defer pgt.Cleanup()

	m.Run()
}

// newTestState is a test helper function that sets up a new State struct with a test
// database and config.
func newTestState(t *testing.T) state.State {
	t.Helper()

	// Note: This will create a fresh database for each test. Currently not an issue,
	// testing the package takes ~1s, but that may change with more tests.
	db, err := pgt.NewDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		t.Fatalf("configure Goose dialect: %v", err)
	}

	if err := goose.Up(db, "../../sql/schema"); err != nil {
		fmt.Fprintf(os.Stderr, "apply migrations: %v\n", err)
	}

	t.Setenv(config.ConfigFilePathEnvVar, t.TempDir()+"/config.json")

	cfg := &config.Config{
		DbURL:           "postgres://test",
		CurrentUserName: "jane",
	}

	return state.State{
		Config: cfg,
		Db:     database.New(db),
	}
}

// createTmpTestConfig creates a new config file in a temp dir and overrides the config
// environment variable with the path to it.
func createTmpTestConfig(t *testing.T) {
	t.Helper()

	// Override the config file with a tempfile
	tmpDir := t.TempDir()
	err := os.Setenv(config.ConfigFilePathEnvVar, tmpDir+"/testconfig.json")
	if err != nil {
		t.Fatalf("unexpected error when setting config env var: %v", err)
	}
}

// TestLogin tests that the login command sets the user to the given argument
func TestLogin(t *testing.T) {
	// Set up the test data
	createTmpTestConfig(t)
	testUser := "john"
	state := newTestState(t)
	commands := SetupRegisterCommands()

	_, err := state.Db.CreateUser(context.Background(), database.CreateUserParams{Name: testUser})
	if err != nil {
		t.Fatalf("unexpected error when creating test user: %v", err)
	}

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
	// Set up the test data
	createTmpTestConfig(t)
	testUser := "alice"
	state := newTestState(t)
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

// TestReset tests that the reset command deletes all users in the database
func TestReset(t *testing.T) {
	// Set up the test data
	createTmpTestConfig(t)
	state := newTestState(t)
	commands := SetupRegisterCommands()

	// Register a random number of users
	for i := range rand.Intn(20) {
		userName := fmt.Sprintf("testUser%d", i)
		cmdRegister := Command{Name: "register", Args: []string{userName}}
		err := commands.Run(&state, cmdRegister)
		if err != nil {
			t.Fatalf("unexpected error when creating test user: %v", err)
		}
	}

	// Run the reset cmdReset
	cmdReset := Command{Name: "reset", Args: []string{}}
	err := commands.Run(&state, cmdReset)
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
	// Set up the test data
	createTmpTestConfig(t)
	testUser := "jill"
	testUserID := uuid.New()
	state := newTestState(t)
	commands := SetupRegisterCommands()

	_, err := state.Db.CreateUser(
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

	// Run the addfeed command
	feedName := "HW Animation Tech Blog"
	feedUrl := "https://hwanimation.tech/feed"
	cmdAddFeed := Command{Name: "addfeed", Args: []string{feedName, feedUrl}}
	err = commands.Run(&state, cmdAddFeed)
	if err != nil {
		t.Fatalf("failed to run the add feed command: %v", err)
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

func TestFeedFollow(t *testing.T) {
	// Set up the test data
	createTmpTestConfig(t)
	userName := "jill"
	userId := uuid.New()
	state := newTestState(t)
	commands := SetupRegisterCommands()

	_, err := state.Db.CreateUser(
		context.Background(),
		database.CreateUserParams{Name: userName, ID: userId},
	)
	if err != nil {
		t.Fatalf("unexpected error when creating test user: %v", err)
	}

	err = state.Config.SetUser(userName)
	if err != nil {
		t.Fatalf("unexpected error when setting current user: %v", err)
	}

	feed, err := state.Db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      "HW Animation Tech Blog",
		Url:       "https://hwanimation.tech/feed",
		UserID:    userId,
	})
	if err != nil {
		t.Fatalf("failed to add feed to database: %v", err)
	}

	// Run the follow command
	cmdFeedFollow := Command{Name: "follow", Args: []string{feed.Url}}
	err = commands.Run(&state, cmdFeedFollow)
	if err != nil {
		t.Fatalf("failed to run the follow command: %v", err)
	}

	// Check that the feed was followed
	got, err := state.Db.GetFeedFollowsForUser(context.Background(), userName)
	if err != nil {
		t.Fatalf("failed to get feed from db: %v", err)
	}
	if len(got) != 1 || got[0].UserName != userName || got[0].FeedName != feed.Name {
		t.Fatalf(`
			"added feed follow data mismatch,
			got %d feed follows, UserName %s, FeedName %s
			want %d feed follows, UserName: %s, FeedName: %s"`,
			len(got), got[0].UserName, got[0].FeedName,
			1, userName, feed.Name)
	}
}
