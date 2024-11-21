package utils

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

// Action interface defines Create and Delete methods
type Action interface {
	Create(source, destination string) error
	Delete(source string) error
}

// SymlinkAction represents the actionable connection between the source and destination
type SymlinkAction struct {
	source      string
	destination string
	isDir       bool
	dryRun      bool
}

func (s SymlinkAction) String(source string) string {
	relSource, _ := filepath.Rel(source, s.source)
	destSource, _ := filepath.Rel(source, s.destination)

	return fmt.Sprintf("%s -> %s", relSource, destSource)
}

// Create creates the symlink.
//
// If dryRun is true, the action will not be executed
func (s SymlinkAction) Create() error {
	return nil
}

// Delete deletes the symlink.
//
// If dryRun is true, the action will not be executed
func (s SymlinkAction) Delete() error {
	return nil
}

// DirAction represents the actions for a directory
type DirAction struct {
	path        string
	permissions fs.FileMode
	dryRun      bool
}

// Create creates the directory.
//
// If dryRun is true, the action will not be executed
func (d DirAction) Create() error {
	return nil
}

// Delete deletes the directory.
//
// If the directory is not empty it will not delete the directory but notify
// the user instead.
//
// If dryRun is true, the action will not be executed
func (d DirAction) Delete() error {
	return nil
}

// Preflight stores actions for the creation of symlinks and directories. It
// also stores any error that occurred during the Preflight so that the user
// can be notified of all errors before failing.
type Preflight struct {
	Errors         []error
	SymlinkActions []SymlinkAction
	DirActions     []DirAction
	DryRun         bool
}

// AddSymlinkAction creates and appends the SymlinkAction to the Preflight check
func (p *Preflight) AddSymlinkAction(source, destination string, isDir bool) {
	p.SymlinkActions = append(p.SymlinkActions, SymlinkAction{
		source:      source,
		destination: destination,
		isDir:       isDir,
		dryRun:      p.DryRun,
	})
}

// AddDirAction creates and appends a DirAction to the Preflight check
func (p *Preflight) AddDirAction(path string, permissions fs.FileMode) {
	p.DirActions = append(p.DirActions, DirAction{
		path:        path,
		permissions: permissions,
		dryRun:      p.DryRun,
	})
}

// AddError appends the provided error to the Preflight check
func (p *Preflight) AddError(err error) {
	p.Errors = append(p.Errors, err)
}
