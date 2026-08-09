package config

import (
	"os"
	"testing"
)

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

func TestGetConfigFilePathEnv(t *testing.T) {
	want := "foo/bar"
	os.Setenv(ConfigFilePathEnvVar, want)

	got, err := getConfigFilePath()
	if got != want || err != nil {
		t.Fatalf("getConfigFilePath() = %q, %v, want match for %#q, nil", got, err, want)
	}
}

// Test a config read/write roundtrip
func TestWriteReadFromFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "tmpConfigDir")
	if err != nil {
		t.Fatalf("Failed to create temp dir during setup: %v", err)
	}

	filePath := tmpDir + "/testconfig.json"
	config := Config{DbURL: "postgres://example", CurrentUserName: "jane"}
	// Override the default config path so we can test using a temp file
	os.Setenv(ConfigFilePathEnvVar, filePath)

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

// Test updating a config user
func TestConfigSetUser(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "tmpConfigDir")
	if err != nil {
		t.Fatalf("Failed to create temp dir during setup: %v", err)
	}

	filePath := tmpDir + "/testconfig.json"
	config := Config{DbURL: "postgres://example", CurrentUserName: "jane"}
	// Override the default config path so we can test using a temp file
	os.Setenv(ConfigFilePathEnvVar, filePath)

	// Write the original config file
	err = writeToFile(filePath, &config)
	if err != nil {
		t.Fatalf("writeToFile(%q, %v) = %v, want match for nil", filePath, config, err)
	}

	// Update the user and check that we get the updated version when reading
	config.SetUser("john")
	got, err := Read()
	if got.CurrentUserName != "john" || err != nil {
		t.Fatalf("readFromFile(%q) = %v, %v, want match for %q, nil", filePath, got, err, config)
	}
}
