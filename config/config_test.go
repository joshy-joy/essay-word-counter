package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const devConfigFilePath = "../resources/dev/config.yml"

const testConfigPath = "./test_config.yml"

// Helper function to create a temporary YAML configuration file
func createTestConfig(content string) error {
	return os.WriteFile(testConfigPath, []byte(content), 0644)
}

// Helper function to clean up the temporary configuration file
func removeTestConfig() {
	os.Remove(testConfigPath)
}

// Test InitConfig with a valid configuration file
func TestInitConfigSuccess(t *testing.T) {
	err := InitConfig(devConfigFilePath)
	assert.Nil(t, err, "Expected no error from InitConfig with valid file")
}

// Test InitConfig with a missing configuration file
func TestInitConfigErrorMissingFile(t *testing.T) {
	err := InitConfig("./non_existent_file.yml")
	assert.NotNil(t, err, "Expected an error due to missing config file")
}

// Test InitConfig with an invalid YAML format
func TestInitConfigErrorInvalidYAML(t *testing.T) {
	invalidYAMLContent := `
							webScrapperJob:
							  count: "invalid_value"  # Invalid, should be an integer
							tokenizerJob:
							  count: 3
							`
	err := createTestConfig(invalidYAMLContent)
	assert.Nil(t, err, "Expected no error while creating the invalid YAML config file")
	err = InitConfig(testConfigPath)
	assert.NotNil(t, err, "Expected an error due to invalid YAML format")
	defer removeTestConfig()
}

// Test Get function to ensure it returns the correct configuration
func TestGetConfigSuccess(t *testing.T) {
	err := InitConfig(devConfigFilePath)
	assert.Nil(t, err, "Expected no error from InitConfig with valid file")
	cfg := Get()
	assert.Equal(t, 2, cfg.WebScrapper.Count, "WebScrapper count should be 2")
	assert.Equal(t, 2, cfg.Tokenizer.Count, "Tokenizer count should be 2")
	assert.Equal(t, int64(30), cfg.External.Timeout, "External timeout should be 30")
	assert.Equal(t, "./example/test.txt", cfg.DefaultFilePath, "Default file path mismatch")
	assert.Equal(t, 2, cfg.ResultLength, "Result length should be 15")
	assert.Equal(t, 3, cfg.WordMinLength, "Word minimum length should be 5")
}

// Test SetFilePath to ensure it updates the DefaultFilePath correctly
func TestSetFilePathSuccess(t *testing.T) {
	err := InitConfig(devConfigFilePath)
	assert.Nil(t, err, "Expected no error from InitConfig with valid file")
	newPath := "./data/new_urls.txt"
	SetFilePath(newPath)
	assert.Equal(t, newPath, Get().DefaultFilePath, "File path should be updated to new path")
}

// Test SetTopN to ensure it updates the ResultLength correctly
func TestSetTopNSuccess(t *testing.T) {
	err := InitConfig(devConfigFilePath)
	assert.Nil(t, err, "Expected no error from InitConfig with valid file")
	newLength := 20
	SetTopN(newLength)
	assert.Equal(t, newLength, Get().ResultLength, "Result length should be updated to new value")
	defer removeTestConfig()
}
