// Package cmd creates the changelog command.
package cmd

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/postfinance/flash"
	"github.com/zbindenren/cc/config"
	"github.com/zbindenren/cc/internal/changelog"
	"github.com/zbindenren/cc/internal/git"
)

const (
	fileOptName           = "f"
	debugOptName          = "d"
	stdOutOptName         = "stdout"
	historyOptName        = "history"
	ignoreOptName         = "ignore"
	sinceTagOptName       = "since"
	initDfltConfigOptName = "init-config"
	noPromptOptName       = "n"
	versionOptName        = "v"
	numOptName            = "num"

	dfltChangelogFile = "CHANGELOG.md"
	dateFormat        = "2006-01-02"
)

// Command represents the changelog CLI.
type Command struct {
	noop bool // for tests
	fs   *flag.FlagSet
	b    BuildInfo
	// flags
	file       *string
	debug      *bool
	toStdOut   *bool
	history    *bool
	ignore     *bool
	sinceTag   *string
	initConfig *bool
	noPrompt   *bool
	version    *bool
	num        *int
}

// New creates a new Command.
func New(b BuildInfo) *Command {
	fs := flag.NewFlagSet("changelog", flag.ExitOnError)

	return &Command{
		fs:         fs,
		b:          b,
		file:       fs.String(fileOptName, dfltChangelogFile, "changelog file name"),
		debug:      fs.Bool(debugOptName, false, "log debug information"),
		toStdOut:   fs.Bool(stdOutOptName, false, "output changelog to stdout instead to file"),
		history:    fs.Bool(historyOptName, false, "create history of old versions tags (output is always stdout)"),
		ignore:     fs.Bool(ignoreOptName, false, "ignore parsing errors of invalid (not conventional) commit messages"),
		sinceTag:   fs.String(sinceTagOptName, "", fmt.Sprintf("in combination with -%s: if a tag is specified, the changelog will be created from that tag on", historyOptName)),
		initConfig: fs.Bool(initDfltConfigOptName, false, fmt.Sprintf("initialize a default changelog configuration '%s'", config.FileName)),
		noPrompt:   fs.Bool(noPromptOptName, false, "do not prompt for next version"),
		version:    fs.Bool(versionOptName, false, "show program version information"),
		num:        fs.Int(numOptName, 0, fmt.Sprintf("in combination with -%s: the number of tags to go back", historyOptName)),
	}
}

// Run parses flags and runs command.
// nolint: gocyclo
func (c Command) Run() error {
	if c.fs != nil {
		if err := c.fs.Parse(os.Args[1:]); err != nil {
			return err
		}
	}

	if *c.version {
		fmt.Println(c.b.Version("changelog"))
		return nil
	}

	// history is always written to stdout
	if *c.history {
		*c.toStdOut = true
	}

	l := flash.New(flash.WithDebug(*c.debug))

	gitCmd, err := git.New(l)
	if err != nil {
		return err
	}

	gitCmd.Noop = c.noop

	if err := c.validate(); err != nil {
		return err
	}

	if *c.initConfig {
		return c.runWriteConfig(l)
	}

	if !gitCmd.IsRepo() {
		return errors.New("current folder is not a git repository")
	}

	cfg, err := config.Load(".")
	if err != nil {
		if err != config.ErrEmpty && err != config.ErrNotFound {
			return err
		}

		l.Debugw("no changelog config file found - using default config", "path", filepath.Join(".", config.FileName))

		cfg = &config.Default
	}

	var dst io.Writer = os.Stdout

	if !*c.toStdOut {
		f, err := os.OpenFile(*c.file, os.O_CREATE|os.O_RDWR, 0600)
		if err != nil {
			return err
		}
		defer f.Close() // nolint: gosec
		dst = f
	}

	if *c.history {
		return c.runHistory(dst, l, *cfg, gitCmd)
	}

	hasTags, err := gitCmd.HasTags()
	if err != nil {
		return err
	}

	if !hasTags {
		return c.runInit(dst, l, *cfg, gitCmd)
	}

	return c.runRelease(dst, l, *cfg, gitCmd)
}

func (c Command) createChangelog(g *git.Command, cfg config.Changelog, l *flash.Logger, revs []string) (*changelog.Changelog, error) {
	cw, err := changelog.New(changelog.WithConfig(cfg), changelog.WithLogFunc(func(msg string, keysAndValues ...interface{}) {
		l.Debugw(msg, keysAndValues...)
	}))
	if err != nil {
		return nil, err
	}

	for _, r := range revs {
		m, err := g.CommitFor(r)
		if err != nil {
			log.Fatal(err)
		}

		if strings.HasPrefix(m.Message, "Merge ") || strings.HasPrefix(m.Message, "Revert ") { // TODO: what else
			continue
		}

		if err := cw.AddMessage(m.Revision, m.Message); err != nil && !*c.ignore {
			return nil, err
		}
	}

	return cw, nil
}

func (c Command) validate() error {
	if *c.sinceTag != "" && !*c.history {
		return fmt.Errorf("'-%s' option is only allowed in combination '-%s' option", sinceTagOptName, historyOptName)
	}

	if *c.num > 0 && !*c.history {
		return fmt.Errorf("'-%s' option is only allowed in combination '-%s' option", numOptName, historyOptName)
	}

	if *c.num > 0 && *c.sinceTag != "" {
		return fmt.Errorf("'-%s' and '-%s' are mutually exclusive", numOptName, sinceTagOptName)
	}

	return nil
}

func (c Command) title(g *git.Command, tag string) (string, error) {
	date, err := g.TagDate(tag)
	if err != nil {
		return "", err
	}

	version, err := semver.NewVersion(tag)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s (%s)", version, date), nil
}

func (c Command) confirmVersion(version semver.Version, in io.Reader, out io.Writer) (*semver.Version, error) {
	if *c.noPrompt {
		return &version, nil
	}

	prompt := fmt.Sprintf("create release %s (press enter to continue with this version or enter version): ", version.String())

	// version prompt, with proposed version
	fmt.Fprint(out, prompt)

	reader := bufio.NewReader(in)

	userInput, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("reading from stdin: %w", err)
	}

	userVersion := strings.TrimSpace(userInput)

	if userVersion == "" {
		return &version, nil
	}

	v, err := semver.NewVersion(userVersion)
	if err != nil {
		return nil, fmt.Errorf("%s is not a valid semantic version: %w", userVersion, err)
	}

	return v, nil
}
