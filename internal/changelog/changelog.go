// Package changelog writes a changelog for conventional commits.
package changelog

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/zbindenren/cc"
	"github.com/zbindenren/cc/config"
)

// ReleaseType determines depending on the commits
// how the version for the next release should be
// increased.
type ReleaseType int

// All possible release types.
const (
	Patch ReleaseType = iota + 1
	Minor
	Major
)

// Changelog creates a changelog.
type Changelog struct {
	githubProjectPath string
	typeSections      typeSections
	cfg               config.Changelog
	releaseType       ReleaseType
	logFunc           func(msg string, keysAndValues ...interface{})
}

// New creates a new Changelog.
func New(opts ...Option) (*Changelog, error) {
	c := Changelog{
		typeSections: typeSections{},
		releaseType:  Patch,
	}

	for _, opt := range opts {
		if err := opt(&c); err != nil {
			return nil, err
		}
	}

	if len(c.cfg.Sections) == 0 {
		c.cfg = config.Default
	}

	return &c, nil
}

// AddMessage add a new commit message to the changelog.
func (c *Changelog) AddMessage(hash, message string) error {
	if c.logFunc != nil {
		c.logFunc("parsing message",
			"ref", hash,
			"msg", message,
		)
	}

	co, err := cc.Parse(message)
	if err != nil {
		return fmt.Errorf("unconventional commit detected - failed to parse '%s': %w", message, err)
	}

	if co.Header.Scope == "" {
		co.Header.Scope = "common"
	}

	if co.Header.Type == "feat" && c.releaseType != Major {
		c.releaseType = Minor
	}

	title, ok := c.cfg.Title(co.Header.Type)
	if !ok {
		title = co.Header.Type
	}

	commit := Commit{
		revisionURL: c.revisionURL(hash),
		Commit:      *co,
	}

	issueNr, ok := closedIssue(*co)
	if ok {
		commit.issueURL = c.issueURL(issueNr)
	}

	if c.logFunc != nil {
		c.logFunc("adding commit",
			"type", co.Header.Type,
			"scope", co.Header.Scope,
			"description", co.Header.Description,
			"body", co.Body,
			"revisonURL", commit.revisionURL,
			"issueURL", commit.issueURL,
		)
	}

	breakingMessage := co.BreakingMessage()
	if breakingMessage != "" {
		c.releaseType = Major
		title = "Breaking Changes"
		c.typeSections.add(title, commit)

		return nil
	}

	c.typeSections.add(title, commit)

	return nil
}

func (c *Changelog) Write(title string, w io.Writer) {
	w.Write([]byte(heading(2, title)))

	for _, title := range c.cfg.List() {
		s, ok := c.typeSections[title]
		if !ok || len(s.scopeSections) == 0 {
			continue
		}

		s.write(w)
	}

	w.Write([]byte(nl + nl + nl))
}

// ReleaseType determines how the version for the next release
// should be increased depending on the added commits.
//
// BREAKING CHANGE: -> Major release
// feat:            -> Minor release
// all other:       -> Patch release
func (c *Changelog) ReleaseType() ReleaseType {
	return c.releaseType
}

// Commit represents a commit.
type Commit struct {
	cc.Commit
	revisionURL string
	issueURL    string
}

type typeSections map[string]typeSection

func (s typeSections) add(title string, c Commit) {
	if _, ok := s[title]; !ok {
		s[title] = typeSection{
			name:          title,
			scopeSections: scopeSections{},
		}
	}

	s[title].scopeSections.add(c)
}

type typeSection struct {
	name          string
	scopeSections scopeSections
}

type scopeSections map[string]*scopeSection

func (s scopeSections) add(c Commit) {
	if _, ok := s[c.Header.Scope]; !ok {
		s[c.Header.Scope] = &scopeSection{
			name:    c.Header.Scope,
			commits: []Commit{},
		}
	}

	s[c.Header.Scope].add(c)
}

func (s scopeSections) list() []scopeSection {
	r := make([]scopeSection, 0, len(s))

	for k := range s {
		r = append(r, *s[k])
	}

	sort.Slice(r, func(i, j int) bool {
		return r[i].name < r[j].name
	})

	return r
}

type scopeSection struct {
	name    string
	commits []Commit
}

func (s *scopeSection) add(c Commit) {
	s.commits = append(s.commits, c)
}

var issueRegexp = regexp.MustCompile(`#(\d+)`)

func closedIssue(c cc.Commit) (issueNR int, ok bool) {
	for _, footer := range c.Footer {
		token := strings.ToLower(footer.Token)
		if strings.Contains(token, "close") || strings.Contains(token, "fix") {
			matches := issueRegexp.FindStringSubmatch(footer.Value)
			if len(matches) < 2 {
				return 0, false
			}

			i, err := strconv.Atoi(matches[1])
			if err != nil {
				return 0, false
			}

			return i, true
		}
	}

	return 0, false
}
