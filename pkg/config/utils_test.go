package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/marcelblijleven/dotlink/pkg/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetTarget_DefaultsToViperTarget(t *testing.T) {
	homedir, err := os.UserHomeDir()
	assert.NoError(t, err)
	viper.SetDefault("target", homedir)
	target := ""

	err = config.GetTarget(&target)

	assert.NoError(t, err)
	assert.Equal(t, homedir, target)
}

func TestGetTarget_ProvidedTargetString(t *testing.T) {
	tempDir := t.TempDir()
	target := tempDir

	err := config.GetTarget(&target)
	assert.NoError(t, err)

	assert.Equal(t, tempDir, target)
}

func TestGetTarget_ProvidedTargetString_DoesNotExist(t *testing.T) {
	tempDir := t.TempDir()
	target := tempDir + "/foo"

	err := config.GetTarget(&target)

	assert.EqualError(t, err, "the provided target directory does not exist")
}

func TestGetTarget_ProvidedTargetString_TargetIsAFile(t *testing.T) {
	tempDir := t.TempDir()
	target := filepath.Join(tempDir, "foo.txt")
	if _, err := os.OpenFile(target, os.O_RDONLY|os.O_CREATE, 0666); err != nil {
		t.Fatalf("error setting up test case: %s", err)
	}

	err := config.GetTarget(&target)

	assert.EqualError(t, err, "the provided target is not a directory")
}

func TestGetTarget_ProvidedTargetString_TargetSameAsCurrent(t *testing.T) {
	current, err := os.Getwd()
	assert.NoError(t, err)

	err = config.GetTarget(&current)

	assert.EqualError(t, err, "the provided target cannot have the current directory as base")
}