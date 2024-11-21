package utils

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/marcelblijleven/dotlink/pkg/config"
)

var FilesMatchErr = errors.New("files match")

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
	compliledIgnorePatterns := CompilePatterns(ignorePatterns)

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

func Scan(source, destination string, ignorePatterns []Pattern, preflight *Preflight) error {
	entries, err := os.ReadDir(source)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourceEntry := filepath.Join(source, entry.Name())
		destinationEntry := filepath.Join(destination, entry.Name())
		// TODO: check for rewrite of path

		if MatchesAnyPattern(sourceEntry, ignorePatterns) {
			continue
		}

		if entry.IsDir() {
			hasIgnoredDescendants, err := HasIgnoredDescendants(sourceEntry, ignorePatterns)
			if err != nil {
				preflight.AddError(fmt.Errorf("error occurred while checking descendants of %s: %v", entry.Name(), err))
				continue
			}

			if !hasIgnoredDescendants {
				// Can create a symlink for the directory
				preflight.AddSymlinkAction(sourceEntry, destinationEntry, true)
			} else {
				// Can't create a symlink, so create the directory manually
				fi, err := os.Stat(sourceEntry)
				if err != nil {
					preflight.AddError(fmt.Errorf("could not retrieve directory permissions for %s: %v", entry.Name(), err))
					continue
				}

				preflight.AddDirAction(destinationEntry, fi.Mode().Perm())

				// Recursively scan directory to see if any of the descendants can be
				// linked efficiently
				Scan(sourceEntry, destinationEntry, ignorePatterns, preflight)
			}
		} else {
			err := checkSymlinkDestination(sourceEntry, destinationEntry)
			if err != nil {
				if !errors.Is(err, FilesMatchErr) {
					preflight.AddError(err)
				}

				continue
			}

			preflight.AddSymlinkAction(sourceEntry, destinationEntry, false)
		}
	}

	return nil
}

// IsSameFile checks if the files are equal, it uses os.Stat to follow symlinks
func IsSameFile(source, destination string) (bool, error) {
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return false, err
	}
	destinationInfo, err := os.Stat(source)
	if err != nil {
		return false, err
	}
	return os.SameFile(sourceInfo, destinationInfo), nil
}

func checkSymlinkDestination(source, destination string) error {
	// Check source
	fi, err := os.Lstat(source)
	if err != nil {
		return fmt.Errorf("cannot create symlink for source %s: %v", fi.Name(), err)
	}

	// Check destination
	_, err = os.Lstat(destination)
	if err == nil {
		same, err := IsSameFile(source, destination)
		if err != nil {
			return fmt.Errorf("cannot create symlink for source %s: %v", err)
		}
		if same {
			return FilesMatchErr
		}
		return fmt.Errorf("cannot create symlink for source %s, destination already exists", source)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("cannot create symlink for source %s: %v", err)
	}

	return nil
}
