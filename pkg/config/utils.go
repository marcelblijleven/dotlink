package config

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ErrTargetDirectoryDoesNotExist       = errors.New("the provided target directory does not exist")
	ErrTargetDirectoryIsNotADirectory    = errors.New("the provided target is not a directory")
	ErrTargetDirectoryIsCurrentDirectory = errors.New("the provided target cannot have the current directory as base")
)

// goosMap maps the value retrieved from goos.GOOS to a common value.
var goosMap = map[string]string{
	"darwin":  "macos",
	"linux":   "linux",
	"windows": "windows",
}

func getMappedGoos() string {
	goos, exists := goosMap[runtime.GOOS]

	if !exists {
		log.Fatalf("'%s' could not be mapped to a supported GOOS value", goos)
	}

	return goos
}

// GetTarget checks if the provided target value is not empty and returns it.
// If the value is empty it will get the target from the configuration, this
// defaults to the user's home directory.
//
// Before returning it will check if the target directory exists and if it is
// a directory.
func GetTarget(target *string) error {
	targetValue := *target

	if targetValue == "" {
		targetValue = viper.GetString("target")
	}

	if strings.HasPrefix(targetValue, "~/") {
		homedir, _ := os.UserHomeDir()
		targetValue = filepath.Join(homedir, targetValue[2:])
	}

	targetValue, err := filepath.Abs(targetValue)
	cobra.CheckErr(err)

	fileinfo, err := os.Stat(targetValue)
	if err != nil {
		return ErrTargetDirectoryDoesNotExist
	}

	if !fileinfo.IsDir() {
		return ErrTargetDirectoryIsNotADirectory
	}

	currentDir, err := os.Getwd()
	cobra.CheckErr(err)

	if strings.Contains(targetValue, currentDir) {
		return ErrTargetDirectoryIsCurrentDirectory
	}

	*target = targetValue
	return nil
}
