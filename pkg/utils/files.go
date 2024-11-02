package utils

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/marcelblijleven/dotlink/pkg/config"
)

type walkFunc = func(root string, walkDirFunc fs.WalkDirFunc) error

// FindFiles walks the provided root directory and finds all files that are not
// ignored by the provided ignore patterns.
func FindFiles(root string, configuration config.Config) ([]string, error) {
	return findFiles(root, configuration.IgnorePatterns(), filepath.WalkDir)
}

// findFiles takes a root to find all files in, ignore patterns to ignore and a
// function that can 'walk' through all the files.
func findFiles(root string, ignorePatterns []string, walker walkFunc) ([]string, error) {
	var files []string
	compliledIgnorePatterns := compilePatterns(ignorePatterns)

	err := walker(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		for _, pattern := range compliledIgnorePatterns {
			if pattern.Matches(path) {
				if d.IsDir() {
					return filepath.SkipDir
				}

				return nil
			}
		}

		if !d.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// IsSymbolicLink checks if the provided path is a symbolic link.
func IsSymbolicLink(path string) (bool, error) {
	stat, err := os.Lstat(path)
	if err != nil {
		return false, err
	}

	return stat.Mode()&fs.ModeSymlink != 0, err
}
