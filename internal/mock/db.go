package mock

import (
	"context"
	"database/sql"
	"fmt"
	"maps"
	"slices"

	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/config"
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/database"
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/state"
	"github.com/google/uuid"
)

type MockDB struct {
	users map[string]database.User
	feeds map[string]database.Feed
}

func (m *MockDB) CreateUser(
	ctx context.Context, arg database.CreateUserParams) (database.User, error) {
	u := database.User(arg)
	m.users[arg.Name] = u
	return u, nil
}

func (m *MockDB) GetUser(ctx context.Context, name string) (database.User, error) {
	u, ok := m.users[name]
	if !ok {
		return database.User{}, sql.ErrNoRows
	}
	return u, nil
}

func (m *MockDB) GetUserByID(ctx context.Context, id uuid.UUID) (database.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return database.User{}, fmt.Errorf("user with id '%v' in database", id)
}

func (m *MockDB) GetUsers(ctx context.Context) ([]database.User, error) {
	users := slices.Collect(maps.Values(m.users))
	return users, nil
}

func (m *MockDB) DeleteAllUsers(ctx context.Context) error {
	m.users = map[string]database.User{}
	return nil
}

func (m *MockDB) CreateFeed(ctx context.Context, arg database.CreateFeedParams) (database.Feed, error) {
	u := database.Feed(arg)
	m.feeds[arg.Url] = u
	return u, nil
}

func (m *MockDB) GetFeed(ctx context.Context, url string) (database.Feed, error) {
	u, ok := m.feeds[url]
	if !ok {
		return database.Feed{}, sql.ErrNoRows
	}
	return u, nil
}

func (m *MockDB) DeleteAllFeeds(ctx context.Context) error {
	m.feeds = map[string]database.Feed{}
	return nil
}

func (m *MockDB) GetFeeds(ctx context.Context) ([]database.Feed, error) {
	feeds := slices.Collect(maps.Values(m.feeds))
	return feeds, nil
}

// GetMockState returns a State object with a mock database, including a feed and a user,
// and a basic config.
func GetMockState() state.State {
	testUser := "jane"
	testFeed := "https://www.example.com/feed"
	cfg := config.Config{DbURL: "postgres://example", CurrentUserName: testUser}
	db := &MockDB{
		users: make(map[string]database.User),
		feeds: make(map[string]database.Feed),
	}
	db.users[testUser] = database.User{Name: testUser}
	db.feeds[testFeed] = database.Feed{Url: testFeed, Name: "Example Feed"}

	return state.State{
		Config: &cfg,
		Db:     db,
	}
}
