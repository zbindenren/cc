// Package config configures the changelog generation.
package config

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

const (
	// FileName is the configuration file name.
	FileName = ".cc.yml"
)

// common errors
var (
	ErrNotFound = errors.New("not found")
	ErrEmpty    = errors.New("empty config")
)

// Changelog configures the changelog.
type Changelog struct {
	Sections          []Section `yaml:"sections"`
	GithubProjectPath string    `yaml:"github_project_path"`
}

// Section is a section config.
type Section struct {
	Type   string `yaml:"type"`
	Title  string `yaml:"title"`
	Hidden bool   `yaml:"hidden"`
}

// Title creates the title from the header type.
func (c Changelog) Title(headerType string) (title string, ok bool) {
	for _, s := range c.Sections {
		if s.Type == headerType {
			return s.Title, true
		}
	}

	return "", false
}

// IsHidden returns true if section should be hidden in changelog.
func (c Changelog) IsHidden(headerType string) bool {
	for _, s := range c.Sections {
		if s.Type == headerType {
			return s.Hidden
		}
	}

	return true
}

// List returns not hidden section titles.
func (c Changelog) List() []string {
	l := make([]string, 0, len(c.Sections))

	for _, s := range c.Sections {
		if !s.Hidden {
			l = append(l, s.Title)
		}
	}

	sort.Strings(l)

	l = append([]string{"Breaking Changes"}, l...)

	return l
}

// Validate validates configuration.
func (c Changelog) Validate() error {
	for _, s := range c.Sections {
		if err := s.validate(); err != nil {
			return err
		}
	}

	return nil
}

// Load is looking for a configuration file named '.cc.yml' in dir. If found
// it tries to unmarshal it into Changelog.
func Load(dir string) (*Changelog, error) {
	configPath := filepath.Join(dir, FileName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, ErrNotFound
	}

	r, err := os.Open(configPath) // nolint: gosec
	if err != nil {
		return nil, fmt.Errorf("open to read %s: %w", configPath, err)
	}

	c, err := Read(r)
	if err != nil {
		return nil, fmt.Errorf("failed to load config from %s: %w", configPath, err)
	}

	if len(c.Sections) == 0 {
		return nil, ErrEmpty
	}

	return c, nil
}

// Read unmarshals config.Changelog from a io.Reader.
func Read(r io.Reader) (*Changelog, error) {
	b, err := ioutil.ReadAll(r) // nolint: gosec
	if err != nil {
		return nil, err
	}

	c := Changelog{}

	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	return &c, nil
}

// Write writes section configration to <dir>/.cc.yml
func Write(dir string, c Changelog) error {
	p := filepath.Join(dir, FileName)

	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(p, b, 0600); err != nil {
		return fmt.Errorf("failed to write file %s: %w", p, err)
	}

	return nil
}

// Default represents the default configuration.
var Default = Changelog{
	Sections: []Section{
		{
			Type:   "build",
			Title:  "Build System",
			Hidden: true,
		},
		{
			Type:   "docs",
			Title:  "Documentation",
			Hidden: true,
		},
		{
			Type:   "feat",
			Title:  "New Features",
			Hidden: false,
		},
		{
			Type:   "fix",
			Title:  "Bug Fixes",
			Hidden: false,
		},
		{
			Type:   "refactor",
			Title:  "Code Refactoring",
			Hidden: true,
		},
		{
			Type:   "test",
			Title:  "Test",
			Hidden: true,
		},
		{
			Type:   "chore",
			Title:  "Tasks",
			Hidden: true,
		},
	},
}

func (s Section) validate() error {
	if s.Title == "" {
		return errors.New("title cannot be empty")
	}

	if s.Type == "" {
		return errors.New("type cannot be empty")
	}

	return nil
}
