package utils

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/marcelblijleven/dotlink/pkg/config"
)

type mapPattern struct {
	Source string
	Target string
}

type CreateLinkOpts interface {
	DryRun() bool
}

// Link represents a link between a source file and target file
type Link struct {
	Source string
	Target string
}

// NewLink creates a Link with the provided source and target
func NewLink(source string, target string) Link {
	return Link{source, target}
}

// String returns a string representation of Link
func (l Link) String() string {
	return fmt.Sprintf("Linking %s => %s", l.Source, l.Target)
}

// Create creates the symbolic link.
func (l Link) Create(dryRun bool) error {
	fmt.Println(l.String())

	if dryRun {
		return nil
	}

	return nil
}

// func CreateLinks(root string, configuration config.Config) ([]Link, error) {
// 	var links []Link
// 	files, err := FindFiles(root, configuration)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	for _, path := range files {
// 		fmt.Println("path", path)
// 		source := filepath.Join(root, path)
// 		target := filepath.Join(configuration.Target, path)
//
// 		mappedTarget, err := mapTarget(target, configuration.MapPatterns())
//
//
// 		links = append(links, NewLink(source, mapTarget(target, configuration.MapPatterns())))
// 	}
//
// 	return links, nil
// }

func mapTarget(target string, mapPatterns []config.MapPattern) (string, error) {
	for _, pattern := range mapPatterns {
		target = strings.Replace(target, pattern.Source, pattern.Target, 1)
	}

	target, err := filepath.Abs(target)
	if err != nil {
		log.Fatalf("could not create mapping for target %s: %e", target, err)
	}

	return target, nil
}
