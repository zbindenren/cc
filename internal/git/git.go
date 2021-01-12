// Package git runs git commands.
package git

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/postfinance/flash"
)

// Command represents a git command.
type Command struct {
	Noop bool // if true no commands are performed that change git state (i.e: no commits and tags are created, no pushs are performed)
	l    *flash.Logger
}

// New Creates a new git command.
func New(l *flash.Logger) (*Command, error) {
	_, err := exec.LookPath("git")
	if err != nil {
		return nil, errors.New("git command not found")
	}

	return &Command{
		l: l,
	}, nil
}

// Run runs the git command.
func (c Command) Run(args ...string) (string, error) {
	extraArgs := []string{
		"-c", "log.showSignature=false",
	}
	args = append(extraArgs, args...)
	cmd := exec.Command("git", args...) // nolint: gosec

	bts, err := cmd.CombinedOutput()
	c.l.Debugw("git command", "args", strings.Join(cmd.Args, " "), "out", string(bts), "err", err)

	if err != nil {
		return "", errors.New(string(bts))
	}

	return string(bts), nil
}

// LastTag returns the last tag.
func (c Command) LastTag(tagMode string) (string, error) {
	if tagMode == "all-branches" {
		tagHash, err := clean(c.Run("rev-list", "--tags", "--max-count=1"))
		if err != nil {
			return "", err
		}

		return clean(c.Run("describe", "--tags", tagHash))
	}

	return clean(c.Run("describe", "--tags", "--abbrev=0"))
}

// IsRepo returns true if current folder is a git repository.
func (c Command) IsRepo() bool {
	out, err := c.Run("rev-parse", "--is-inside-work-tree")
	return err == nil && strings.TrimSpace(out) == "true"
}

// TagDate gets the date of the tag.
func (c Command) TagDate(tag string) (string, error) {
	if tag == "" {
		return "", errors.New("tag cannot be empty")
	}

	out, err := clean(c.Run("log", "-1", "--format=%ai", tag))
	if err != nil {
		return "", err
	}

	return strings.Split(out, " ")[0], nil
}

// HasUncommitted checks if there are uncommitted changes.
func (c Command) HasUncommitted() (bool, error) {
	out, err := c.Run("diff-index", "HEAD")
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(out) != "", nil
}

// HasRemotes checks if a remote is configured.
func (c Command) HasRemotes() bool {
	_, err := c.Run("ls-remote", "--exit-code")
	return err == nil
}

// IsStaged checks if a path is staged in repository.
func (c Command) IsStaged(path string) bool {
	_, err := c.Run("ls-files", "--error-unmatch", path)
	return err == nil
}

// RevList runs git rev-list start..end
func (c Command) RevList(start, end string) ([]string, error) {
	arg := end
	if start != "" {
		arg = start + ".." + end
	}

	revs, err := c.Run("rev-list", arg)
	if err != nil {
		return nil, err
	}

	revs = strings.TrimSpace(revs)

	if revs == "" {
		return []string{}, nil
	}

	return strings.Split(strings.TrimSpace(revs), "\n"), nil
}

// CreateRelease creates a release tag.
func (c Command) CreateRelease(version string) error {
	cmd := []string{"git", "tag", "-a", "v" + version, "-m", fmt.Sprintf("chore: bump version to %s", version)}

	if c.Noop {
		c.l.Debugw("noop mode - command not run", "cmd", strings.Join(cmd, " "))
		return nil
	}

	_, err := c.Run(cmd[1:]...)

	return err
}

// CommitFile commits a file.
func (c Command) CommitFile(file, msg string) error {
	cmd := []string{"git", "commit", file, "-m", msg}

	if c.Noop {
		c.l.Debugw("noop mode - command not run", "cmd", strings.Join(cmd, " "))
		return nil
	}

	_, err := c.Run(cmd[1:]...)

	return err
}

// StageFile stages a file.
func (c Command) StageFile(file string) error {
	cmd := []string{"git", "add", file}

	if c.Noop {
		c.l.Debugw("noop mode - command not run", "cmd", strings.Join(cmd, " "))
		return nil
	}

	_, err := c.Run(cmd[1:]...)

	return err
}

// Push pushes tags and commits.
func (c Command) Push() error {
	cmd := []string{"git", "push", "--follow-tags"}

	if c.Noop {
		c.l.Debugw("noop mode - command not run", "cmd", strings.Join(cmd, " "))
		return nil
	}

	_, err := c.Run(cmd[1:]...)

	return err
}

// ListTags list tags.
func (c Command) ListTags() (Tags, error) {
	tags, err := c.Run("tag", "--sort=-v:refname")
	if err != nil {
		return nil, err
	}

	return strings.Split(strings.TrimSpace(tags), "\n"), nil
}

// HasTags returns true if repository is tagged.
func (c Command) HasTags() (bool, error) {
	out, err := c.Run("tag")
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(out) != "", nil
}

// CommitFor creates a Commit for a revision.
func (c Command) CommitFor(revision string) (*Commit, error) {
	m, err := c.Run("show", "--format=%B", "-s", revision)
	if err != nil {
		return nil, fmt.Errorf("failed to get message for '%s': %w", revision, err)
	}

	return &Commit{
		Revision: revision,
		Message:  strings.TrimSpace(m),
	}, nil
}

// Tags is a slice of tags.
type Tags []string

// Index returns the index of name. If not found -1 is returned.
func (t Tags) Index(name string) int {
	for i := range t {
		if t[i] == name {
			return i
		}
	}

	return -1
}

// Commit represents a commit (revision and message).
type Commit struct {
	Message  string
	Revision string
}

func clean(output string, err error) (string, error) {
	output = strings.ReplaceAll(strings.Split(output, "\n")[0], "'", "")

	if err != nil {
		err = errors.New(strings.TrimSuffix(err.Error(), "\n"))
	}

	return output, err
}
