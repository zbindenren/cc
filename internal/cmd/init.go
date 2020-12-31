package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Masterminds/semver"
	"github.com/postfinance/flash"
	"github.com/zbindenren/cc/config"
	"github.com/zbindenren/cc/internal/git"
)

func (c Command) runInit(dst io.Writer, l *flash.Logger, cfg config.Changelog, g *git.Command) error {
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

	next, _ := semver.NewVersion("v0.1.0")

	revs, err := g.RevList("", "HEAD")
	if err != nil {
		return err
	}

	cw, err := c.createChangelog(g, cfg, l, revs)
	if err != nil {
		return err
	}

	fmt.Printf("create first version: %s\n", next)

	version, err := c.confirmVersion(*next, os.Stdin, os.Stdout)
	if err != nil {
		return err
	}

	title := fmt.Sprintf("%s (%s)", version, time.Now().Format(dateFormat))

	cw.Write(title, dst)

	if !*c.toStdOut {
		l.Debugw("staging file", "file", *c.file)

		if err := g.StageFile(*c.file); err != nil {
			return err
		}

		l.Debug("committing changes")

		if err := g.CommitFile(*c.file, fmt.Sprintf("chore: update changelog with %s release", next.String())); err != nil {
			return err
		}

		l.Debugw("create release tag", "release", "v"+next.String())

		if err := g.CreateRelease(next.String()); err != nil {
			return err
		}

		l.Debug("pushing tags")

		if err := g.Push(); err != nil {
			return err
		}
	}

	return nil
}
