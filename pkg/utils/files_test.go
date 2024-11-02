package utils

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockWalkDir(paths []string, errPaths map[string]error) walkFunc {
	return func(root string, walkDirFunc fs.WalkDirFunc) error {
		for _, path := range paths {
			var d fs.DirEntry
			var err error

			// Mock directory entries
			if errPaths != nil && errPaths[path] != nil {
				err = errPaths[path]
			} else {
				d = newMockDirEntry(path)
			}

			walkErr := walkDirFunc(path, d, err)
			if walkErr == filepath.SkipDir {
				break
			} else if walkErr != nil {
				return walkErr
			}
		}
		return nil
	}
}

// Mock directory entry structure
type mockDirEntry struct {
	name  string
	isDir bool
}

func (m mockDirEntry) Name() string               { return m.name }
func (m mockDirEntry) IsDir() bool                { return m.isDir }
func (m mockDirEntry) Type() fs.FileMode          { return 0 }
func (m mockDirEntry) Info() (fs.FileInfo, error) { return nil, nil }

func newMockDirEntry(path string) fs.DirEntry {
	return mockDirEntry{name: path, isDir: filepath.Ext(path) == ""}
}

func TestFindFiles_BasicFunctionality(t *testing.T) {
	paths := []string{"file1.txt", "file2.txt", "dir/file3.txt"}
	walker := mockWalkDir(paths, nil)

	files, err := findFiles("/", []string{}, walker)
	assert.NoError(t, err)

	expectedFiles := []string{"file1.txt", "file2.txt", "dir/file3.txt"}
	assert.Equal(t, expectedFiles, files)
}

func TestFindFiles_IgnoreRegexPattern(t *testing.T) {
	paths := []string{"file1.txt", "file2.log", "dir/file3.log"}
	walker := mockWalkDir(paths, nil)

	files, err := findFiles("/", []string{`\.log$`}, walker)
	assert.NoError(t, err)

	expectedFiles := []string{"file1.txt"}
	assert.Equal(t, expectedFiles, files)
}

func TestFindFiles_IgnoreLiteralPattern(t *testing.T) {
	paths := []string{"file1.txt", "node_modules/file2.txt", "file3.txt"}
	walker := mockWalkDir(paths, nil)

	files, err := findFiles("/", []string{"node_modules"}, walker)
	assert.NoError(t, err)

	expectedFiles := []string{"file1.txt", "file3.txt"}
	assert.Equal(t, expectedFiles, files)
}

func TestFindFiles_EmptyDirectory(t *testing.T) {
	walker := mockWalkDir([]string{}, nil)

	files, err := findFiles("/", []string{}, walker)

	assert.NoError(t, err)
	assert.Len(t, files, 0)
}

func TestFindFiles_ErrorHandling(t *testing.T) {
	paths := []string{"file1.txt", "dir/file2.txt"}
	errPaths := map[string]error{"dir/file2.txt": errors.New("permission denied")}
	walker := mockWalkDir(paths, errPaths)

	_, err := findFiles("/", []string{}, walker)

	assert.Error(t, err)
}

func TestIsSymbolicLink(t *testing.T) {
	tempDirA := t.TempDir()
	tempDirB := t.TempDir()
	file := filepath.Join(tempDirA, "no-link.txt")
	link := filepath.Join(tempDirB, "link.txt")
	f, err := os.OpenFile(file, os.O_CREATE, 0666)
	assert.NoError(t, err)
	f.Close()

	check, err := IsSymbolicLink(file)
	assert.NoError(t, err)
	assert.False(t, check)
	err = os.Symlink(file, link)
	assert.NoError(t, err)

	check, err = IsSymbolicLink(link)
	assert.NoError(t, err)
	assert.True(t, check)
}
