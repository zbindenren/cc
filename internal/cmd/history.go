package cmd

import (
	"fmt"
	"io"

	"github.com/postfinance/flash"
	"github.com/zbindenren/cc/config"
	"github.com/zbindenren/cc/internal/git"
)

func (c Command) runHistory(dst io.Writer, l *flash.Logger, cfg config.Changelog, g *git.Command) error {
	tags, err := g.ListTags()
	if err != nil {
		return err
	}

	max := len(tags) - 1
	if *c.sinceTag != "" {
		max = tags.Index(*c.sinceTag)

		if max < 0 {
			return fmt.Errorf("tag '%s' not found", *c.sinceTag)
		}
	}

	for i := 0; i < max; i++ {
		revs, err := g.RevList(tags[i+1], tags[i])
		if err != nil {
			return err
		}

		cw, err := c.createChangelog(g, cfg, l, revs)
		if err != nil {
			return err
		}

		title, err := c.title(g, tags[i])
		if err != nil {
			return err
		}

		cw.Write(title, dst)
	}

	return nil
}
