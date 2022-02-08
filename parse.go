// Package cc parses conventional commits: https://www.conventionalcommits.org/en/v1.0.0
//
// A commit is of the form:
//	   <type>[optional scope]: <description>
//
//	   [optional body]
//
//	   [optional footer(s)]
//
// Example:
//	   fix(cluster): connection resets
//
//	   here an additional body can be added.
//	   the body can be multiline.
//
//	   Acknowledged-by: a user
//	   Fixes #123 bug
package cc

import (
	"bytes"
	"strings"

	"github.com/bbuck/go-lexer"
)

// Commit contains the parsed conventional commit data.
type Commit struct {
	Header     Header
	Body       string
	Footer     Footers
	isBreaking bool
}

// Header is the header part of the conventional commit.
type Header struct {
	Type        string
	Scope       string
	Description string
}

// Footer is the footer (trailer) part of the conventional commit.
type Footer struct {
	Token string `yaml:"token"`
	Value string `yaml:"value"`
}

// Footers is a slice of Footer.
type Footers []Footer

// BreakingMessage returns the breaking change message text. If the breaking change
// is indicated with a `BREAKING CHANGE` or `BREAKING_CHANGE` footer token, then this
// value is returned. Otherwise the header description is returned.
//
// If no breaking change is detected an empty string is returned.
func (c Commit) BreakingMessage() string {
	b := c.Footer.breakingMessage()
	if b == "" && c.isBreaking {
		b = c.Header.Description
	}

	return b
}

// Parse parses the conventional commit. If it fails, an error is returned.
func Parse(s string) (*Commit, error) {
	l := lexer.New(normalizeNewlines(strings.TrimSpace(s)), typeState)
	l.ErrorHandler = func(string) {}

	l.Start()

	c := Commit{}
	footerCount := 0

	for {
		t, done := l.NextToken()
		if done {
			break
		}

		switch t.Type {
		case breakingChange:
			c.isBreaking = true
		case headerScope:
			c.Header.Scope = t.Value
		case headerType:
			c.Header.Type = t.Value
		case description:
			c.Header.Description = t.Value
		case body:
			c.Body = strings.TrimSpace(t.Value)
		case footerToken:
			c.Footer = append(c.Footer, Footer{
				Token: t.Value,
			})
		case footerValue:
			c.Footer[footerCount].Value = strings.TrimSpace(t.Value)
			footerCount++
		}
	}

	if l.Err != nil {
		return nil, l.Err
	}

	return &c, nil
}

func (f Footer) isBreaking() bool {
	return strings.Replace(f.Token, "-", " ", 1) == breakingChangeFooterToken
}

func (f Footers) breakingMessage() string {
	for _, footer := range f {
		if footer.isBreaking() {
			return footer.Value
		}
	}

	return ""
}

func normalizeNewlines(s string) string {
	d := []byte(s)
	// replace CR LF \r\n (windows) with LF \n (unix)
	d = bytes.ReplaceAll(d, []byte{13, 10}, []byte{10})
	// replace CF \r (mac) with LF \n (unix)
	d = bytes.ReplaceAll(d, []byte{13}, []byte{10})

	return string(d)
}
