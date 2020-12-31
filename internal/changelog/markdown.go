package changelog

import (
	"bufio"
	"fmt"
	"io"
	"path"
	"sort"
	"strings"
)

const (
	tabStop    = "  "
	boldPrefix = "**"
	bullet     = '*'
	nl         = '\n'

	githubURL        = "https://github.com"
	githubCommitPath = "/commit"
	githubIssuesPath = "/issues"
)

func (s typeSection) write(w io.Writer) {
	var sb strings.Builder

	sb.WriteString("\n\n")
	sb.WriteString(heading(3, s.name))
	sb.WriteRune(nl)

	w.Write([]byte(sb.String()))

	for _, scope := range s.scopeSections.list() {
		scope.write(w)
	}
}

func (s scopeSection) write(w io.Writer) {
	var b strings.Builder

	commits := s.commits
	sort.Slice(commits, func(i, j int) bool {
		if commits[i].Header.Scope == commits[j].Header.Scope {
			return commits[i].Header.Description < commits[j].Header.Description
		}

		return commits[i].Header.Scope < commits[j].Header.Scope
	})

	for i := range commits {
		if i == 0 && commits[i].BreakingMessage() != "" {
			b.WriteString(listItem(1, bold(s.name)))
			w.Write([]byte(b.String()))
		}

		commits[i].write(w)
	}
}

func (c *Commit) write(w io.Writer) {
	if c.BreakingMessage() != "" {
		c.writeBreaking(w)
		return
	}

	c.writeStandard(w)
}

func (c *Commit) writeBreaking(w io.Writer) {
	var s strings.Builder

	s.WriteString(bold(c.revisionURL))
	s.WriteRune(':')
	s.WriteRune(nl)
	s.WriteString(c.Header.Description)

	if c.issueURL != "" {
		s.WriteString(" (")
		s.WriteString(c.issueURL)
		s.WriteString(")")
	}

	if c.Body != "" {
		s.WriteRune(nl)
		s.WriteString(blockQuote(0, c.Body))
	}

	l := listItem(2, s.String())
	w.Write([]byte(l))
}

func (c *Commit) writeStandard(w io.Writer) {
	var s strings.Builder

	urls := []string{c.revisionURL}
	if c.issueURL != "" {
		urls = append([]string{c.issueURL}, urls...)
	}

	s.WriteString(bold(c.Header.Scope))
	s.WriteString(": ")
	s.WriteString(c.Header.Description)
	s.WriteString(" (")
	s.WriteString(strings.Join(urls, ", "))
	s.WriteString(")")

	if c.Body != "" {
		s.WriteRune(nl)
		s.WriteString(blockQuote(0, c.Body))
	}

	l := listItem(1, s.String())
	w.Write([]byte(l))
}

func (c Changelog) issueURL(issueNumber int) string {
	// gitlab renders link automatically for revisions
	if c.githubProjectPath == "" {
		return fmt.Sprintf("#%d", issueNumber)
	}

	var (
		s          strings.Builder
		issuesPath = githubIssuesPath
	)

	url := path.Join(githubURL, c.githubProjectPath, issuesPath) + "/"
	issue := fmt.Sprintf("#%d", issueNumber)

	s.WriteRune('[')
	s.WriteString(issue)
	s.WriteRune(']')
	s.WriteRune('(')
	s.WriteString(url)
	s.WriteString(issue)
	s.WriteRune(')')

	return s.String()
}
func (c Changelog) revisionURL(revision string) string {
	// gitlab renders link automatically for revisions
	if c.githubProjectPath == "" {
		return revision
	}

	var (
		s          strings.Builder
		commitPath = githubCommitPath
	)

	url := path.Join(githubURL, c.githubProjectPath, commitPath) + "/"

	s.WriteRune('[')
	s.WriteString(revision[:8])
	s.WriteRune(']')
	s.WriteRune('(')
	s.WriteString(url)
	s.WriteString(revision)
	s.WriteRune(')')

	return s.String()
}

func heading(level int, title string) string {
	var s strings.Builder

	s.WriteString(strings.Repeat("#", level))
	s.WriteRune(' ')
	s.WriteString(title)
	s.WriteRune('\n')

	return s.String()
}

func listItem(level int, data string) string {
	var s strings.Builder

	s.WriteString(strings.Repeat(tabStop, level-1))
	s.WriteRune(bullet)
	s.WriteString(" ")

	scanner := bufio.NewScanner(strings.NewReader(data))
	ident := strings.Repeat(tabStop, level-1) + "  "
	count := 0

	for scanner.Scan() {
		if count > 0 {
			s.WriteString(ident)
		}

		s.WriteString(scanner.Text())
		s.WriteRune(nl)
		count++
	}

	return s.String()
}

func bold(data string) string {
	var s strings.Builder

	s.WriteString(boldPrefix)
	s.WriteString(data)
	s.WriteString(boldPrefix)

	return s.String()
}

func blockQuote(indent int, data string) string {
	var s strings.Builder

	scanner := bufio.NewScanner(strings.NewReader(data))
	ident := strings.Repeat(tabStop, indent)

	for scanner.Scan() {
		s.WriteString(ident)
		s.WriteString("> ")
		s.WriteString(scanner.Text())
		s.WriteRune(nl)
	}

	return s.String()
}
