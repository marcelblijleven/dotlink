package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/marcelblijleven/dotlink/pkg/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setupConfigFile(t *testing.T) string {
	tempDir := t.TempDir()
	viper.AddConfigPath(tempDir)
	os.WriteFile(filepath.Join(tempDir, ".dotlink.yaml"), []byte(fmt.Sprintf("target: %s", tempDir)), 0755)
	return tempDir
}

func TestInitConfig_NoConfigFileProvided_UsesDefaultFile(t *testing.T) {
	tempDir := setupConfigFile(t)
	target := ""
	configuration := config.Config{}

	config.InitConfig("", &target, nil, &configuration)

	assert.Equal(t, configuration.Target, tempDir)
	assert.Equal(t, configuration.IgnorePatterns(), []string{".git", ".dotlink.yaml"})
}

func TestInitConfig_ConfigFileProvided(t *testing.T) {
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, ".provided_dotlink.yaml")
	os.WriteFile(file, []byte(fmt.Sprintf("target: %s", tempDir)), 0755)

	target := ""
	configuration := config.Config{}

	config.InitConfig(file, &target, nil, &configuration)

	assert.Equal(t, configuration.Target, tempDir)
	assert.Equal(t, configuration.IgnorePatterns(), []string{".git", ".dotlink.yaml"})
}

func TestInitConfig_IgnorePatterns_AreAdded(t *testing.T) {
	setupConfigFile(t)
	target := ""
	configuration := config.Config{}
	ignorePatterns := []string{"foo", "bar"}

	config.InitConfig("", &target, ignorePatterns, &configuration)

	assert.Equal(t, configuration.IgnorePatterns(), []string{".git", ".dotlink.yaml", "foo", "bar"})
}
