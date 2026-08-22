package mock

import (
	"context"
	"database/sql"
	"maps"
	"slices"

	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/config"
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/database"
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/state"
)

type MockDB struct {
	users map[string]database.User
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

func (m *MockDB) GetUsers(ctx context.Context) ([]database.User, error) {
	users := slices.Collect(maps.Values(m.users))
	return users, nil
}

func (m *MockDB) DeleteAllUsers(ctx context.Context) error {
	m.users = map[string]database.User{}
	return nil
}

func (m *MockDB) CreateFeed(ctx context.Context, arg database.CreateFeedParams) (database.Feed, error) {
	return database.Feed{}, nil
}

func (m *MockDB) GetFeed(ctx context.Context, name string) (database.Feed, error) {
	return database.Feed{}, nil
}

func (m *MockDB) DeleteAllFeeds(ctx context.Context) error {
	return nil
}

func (m *MockDB) GetFeeds(ctx context.Context) ([]database.Feed, error) {
	return []database.Feed{}, nil
}

func GetMockState() state.State {
	testUser := "jane"
	cfg := config.Config{DbURL: "postgres://example", CurrentUserName: testUser}
	db := &MockDB{
		users: make(map[string]database.User),
	}
	db.users[testUser] = database.User{Name: testUser}

	return state.State{
		Config: &cfg,
		Db:     db,
	}
}
