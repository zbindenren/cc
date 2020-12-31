package changelog

import (
	"errors"

	"github.com/zbindenren/cc/config"
)

// Option is a functional option.
type Option func(*Changelog) error

// LogFunc is a logging function.
type LogFunc func(msg string, keysAndValues ...interface{})

// WithGithubProjectPath configures the project path for github. This determines how links to issues and
// commits are created.
func WithGithubProjectPath(projectPath string) Option {
	return func(c *Changelog) error {
		if projectPath == "" {
			return errors.New("project path is empty")
		}

		c.githubProjectPath = projectPath

		return nil
	}
}

// WithConfig configures how header types are mapped to headings and which
// sections are hidden.
func WithConfig(cfg config.Changelog) Option {
	return func(c *Changelog) error {
		c.cfg = cfg
		return nil
	}
}

// WithLogFunc can be used to configure a logging function to show debug output
// when adding a git message.
func WithLogFunc(f LogFunc) Option {
	return func(c *Changelog) error {
		if f == nil {
			return errors.New("log func cannot be nil")
		}

		c.logFunc = f

		return nil
	}
}
