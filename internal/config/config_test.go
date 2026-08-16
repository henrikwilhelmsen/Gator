package config

import (
	"os"
	"testing"
)

// TestGetConfigFilePath tests that the getConfigFilePath function returns the expected path
func TestGetConfigFilePath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get UserHomeDir during setup: %v", err)
	}

	want := homeDir + "/.gatorconfig.json"
	got, err := getConfigFilePath()
	if got != want || err != nil {
		t.Fatalf("getConfigFilePath() = %q, %v, want match for %#q, nil", got, err, want)
	}
}

// TestGetConfigFilePathEnv tests that we can override the config
// path with an environment variable
func TestGetConfigFilePathEnv(t *testing.T) {
	want := "foo/bar"
	err := os.Setenv(ConfigFilePathEnvVar, want)
	if err != nil {
		t.Fatalf("unexpected error when setting config env var: %v", err)
	}

	got, err := getConfigFilePath()
	if got != want || err != nil {
		t.Fatalf("getConfigFilePath() = %q, %v, want match for %#q, nil", got, err, want)
	}
}

// TestWriteReadFromFile tests a config read/write roundtrip
func TestWriteReadFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := tmpDir + "/testconfig.json"
	config := Config{DbURL: "postgres://example", CurrentUserName: "jane"}
	// Override the default config path so we can test using a temp file
	err := os.Setenv(ConfigFilePathEnvVar, filePath)
	if err != nil {
		t.Fatalf("unexpected error when setting config env var: %v", err)
	}

	// Test that we can write to file with no errors
	err = writeToFile(filePath, &config)
	if err != nil {
		t.Fatalf("writeToFile(%q, %v) = %v, want match for nil", filePath, config, err)
	}

	// Test that we can read the file and get the same config back
	got, err := Read()
	if got != config || err != nil {
		t.Fatalf("readFromFile(%q) = %v, %v, want match for %q, nil", filePath, got, err, config)
	}
}

// TestConfigSetUser tests that the SetUser method updates the config file and that the
// Read function returns an updated config.
func TestConfigSetUser(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := tmpDir + "/testconfig.json"
	config := Config{DbURL: "postgres://example", CurrentUserName: "jane"}
	// Override the default config path so we can test using a temp file
	err := os.Setenv(ConfigFilePathEnvVar, filePath)
	if err != nil {
		t.Fatalf("unexpected error when setting config env var: %v", err)
	}

	// Write the original config file
	err = writeToFile(filePath, &config)
	if err != nil {
		t.Fatalf("writeToFile(%q, %v) = %v, want match for nil", filePath, config, err)
	}

	// Update the user and check that we get the updated version when reading
	err = config.SetUser("john")
	if err != nil {
		t.Fatalf("unexpected error when setting config user: %v", err)
	}

	got, err := Read()
	if got.CurrentUserName != "john" || err != nil {
		t.Fatalf("readFromFile(%q) = %v, %v, want match for %q, nil", filePath, got, err, config)
	}
}
