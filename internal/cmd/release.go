package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/Masterminds/semver"
	"github.com/postfinance/flash"
	"github.com/zbindenren/cc/config"
	"github.com/zbindenren/cc/internal/changelog"
	"github.com/zbindenren/cc/internal/git"
)

// nolint: gocyclo,funlen
func (c Command) runRelease(dst io.Writer, l *flash.Logger, cfg config.Changelog, g *git.Command) error {
	if !g.HasRemotes() {
		return errors.New("git repo has no remotes configured, cannot initialize changelog")
	}

	uncommmited, err := g.HasUncommitted()
	if err != nil {
		return err
	}

	if uncommmited {
		return errors.New("git repository contains uncommitted changes")
	}

	var next semver.Version

	tag, err := g.LastTag("current-branch")
	if err != nil {
		return err
	}

	revs, err := g.RevList("tags/"+tag, "HEAD")
	if err != nil {
		return err
	}

	cw, err := c.createChangelog(g, cfg, l, revs)
	if err != nil {
		return err
	}

	current, err := semver.NewVersion(tag)
	if err != nil {
		return err
	}

	switch cw.ReleaseType() {
	case changelog.Patch:
		next = current.IncPatch()
	case changelog.Minor:
		next = current.IncMinor()
	case changelog.Major:
		next = current.IncMajor()
	}

	fmt.Printf("last version: %s\n", current)
	fmt.Printf("next version: %s\n", &next)

	version, err := c.confirmVersion(next, os.Stdin, os.Stdout)
	if err != nil {
		return err
	}

	if !version.GreaterThan(current) {
		return fmt.Errorf("version must be greater than current version %s", current)
	}

	title := fmt.Sprintf("%s (%s)", version, time.Now().Format(dateFormat))

	old, err := ioutil.ReadFile(*c.file)
	if err != nil && !*c.toStdOut {
		return err
	}

	l.Debugw("update changelog", "file", *c.file, "title", title)
	cw.Write(title, dst)

	if !*c.toStdOut {
		if _, err := dst.Write(old); err != nil {
			return err
		}

		if !g.IsStaged(*c.file) {
			l.Debug("staging changelog", "file", *c.file)

			if err := g.StageFile(*c.file); err != nil {
				return err
			}
		}

		l.Debug("committing changes")

		if err := g.CommitFile(*c.file, fmt.Sprintf("chore: update changelog with %s release", version.String())); err != nil {
			return err
		}

		l.Debugw("create release tag", "release", "v"+version.String())

		if err := g.CreateRelease(version.String()); err != nil {
			return err
		}

		l.Debug("pushing tags")

		if err := g.Push(); err != nil {
			return err
		}
	}

	return nil
}
