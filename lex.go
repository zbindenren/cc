package cc

import (
	"fmt"
	"unicode"

	"github.com/bbuck/go-lexer"
)

const (
	breakingChangeFooterToken = "BREAKING CHANGE"
)

// all lexer token types that are emitted.
const (
	body lexer.TokenType = iota
	breakingChange
	descriptionDelimiter
	description
	footerDelimter
	footerToken
	footerValue
	leftScopeDelimiter
	rightScopeDelimiter
	headerScope
	headerType
)

func typeState(l *lexer.L) lexer.StateFunc {
	for {
		r := l.Peek()

		if r == lexer.EOFRune {
			l.Error("missing scope or description")

			return nil
		}

		if r == ':' {
			l.Emit(headerType)

			return descriptionDelimiterState
		}

		if r == '!' {
			l.Emit(headerType)

			return descriptionDelimiterState
		}

		if r == '(' {
			l.Emit(headerType)
			l.Take("(")
			l.Emit(leftScopeDelimiter)

			return scopeState
		}

		if !unicode.IsLetter(r) {
			l.Error(fmt.Sprintf("invalid character '%c' in type", r))

			return nil
		}

		l.Next()
	}
}

func scopeState(l *lexer.L) lexer.StateFunc {
	for {
		r := l.Peek()
		if r == ')' {
			if l.Current() == "" {
				l.Error("empty scope")

				return nil
			}

			l.Emit(headerScope)
			l.Take(")")
			l.Emit(rightScopeDelimiter)

			return descriptionDelimiterState
		}

		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_') {
			l.Error("scope must be noun in ()")

			return nil
		}

		l.Next()
	}
}

func descriptionDelimiterState(l *lexer.L) lexer.StateFunc {
	if l.Peek() == '!' {
		l.Next()
		l.Emit(breakingChange)
	}

	l.Take(": ")

	if l.Current() != ": " {
		l.Error("scope must be followed by ': '")
		return nil
	}

	l.Emit(descriptionDelimiter)

	return descriptionState
}

func descriptionState(l *lexer.L) lexer.StateFunc {
	for {
		if l.Next() == lexer.EOFRune {
			l.Emit(description)
			return nil
		}

		if l.Peek() == '\n' {
			l.Emit(description)
			return headerDelimeterState
		}
	}
}

func headerDelimeterState(l *lexer.L) lexer.StateFunc {
	l.Take("\n")

	if len(l.Current()) < 2 {
		l.Error("at least one empty line required after header")
		return nil
	}

	l.Ignore()

	return bodyOrFooterState
}

func bodyOrFooterState(l *lexer.L) lexer.StateFunc {
	count := takeFooterToken(l)

	// there is no body
	if count > 0 {
		rewind(l, count)
		return footerTokenState
	}

	return bodyState
}

func bodyState(l *lexer.L) lexer.StateFunc {
	found := takeUntilFirstFooterToken(l)
	if !found {
		l.Emit(body)
		return nil
	}

	// go back to the last newline character
	for {
		l.Rewind()

		if l.Peek() != '\n' {
			break
		}
	}
	l.Next()
	l.Emit(body)

	return bodyDelimeterState
}

func bodyDelimeterState(l *lexer.L) lexer.StateFunc {
	l.Take("\n")

	if len(l.Current()) < 2 {
		l.Error("at least one empty line required after body")
		return nil
	}

	l.Ignore()

	return footerTokenState
}

func footerTokenState(l *lexer.L) lexer.StateFunc {
	l.Take("\n")
	l.Ignore()

	takeFooterToken(l)
	l.Emit(footerToken)

	return footerDelimiterState
}

func footerValueState(l *lexer.L) lexer.StateFunc {
	if l.Peek() == lexer.EOFRune {
		return nil
	}

	found := takeUntilFirstFooterToken(l)
	l.Emit(footerValue)

	if !found {
		return nil
	}

	return footerTokenState
}

func footerDelimiterState(l *lexer.L) lexer.StateFunc {
	l.Take(": ")
	l.Emit(footerDelimter)

	return footerValueState
}

// takeUntilFirstFooter takes all characters until a footer token is detected
func takeUntilFirstFooterToken(l *lexer.L) bool {
	for {
		r := l.Next()
		if r == lexer.EOFRune {
			return false
		}

		// a footer token has to begin at the start of a line
		if r == '\n' {
			count := takeFooterToken(l)
			if count > 0 {
				// if count is > 0 we are at the end of the footer token
				for i := 0; i < count; i++ {
					l.Rewind()
				}

				return true
			}
		}
	}
}

func rewind(l *lexer.L, count int) {
	for i := count; i > 0; i-- {
		l.Rewind()
	}
}

// takeFooterToken continues over each consecutive character in the source
// until an invalid footer token character is detected. The method returns
// the length if there is a valid footer token found. If it is not a valid footer
// token, 0 is returned.
func takeFooterToken(l *lexer.L) int {
	// handle BREAKING CHANGE: (BREAKING-CHANGE: is handled as standard footer token)
	if l.Peek() == 'B' {
		candidate := peekString(l, len(breakingChangeFooterToken)+2)
		if candidate == breakingChangeFooterToken+": " {
			l.Take(breakingChangeFooterToken)

			return len(breakingChangeFooterToken)
		}
	}

	count := 0
	r := l.Next()

	for unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
		count++

		r = l.Next()
	}

	// :<space> or <space># delimiter
	if (r == ':' && l.Peek() == ' ') || (r == ' ' && l.Peek() == '#') {
		l.Rewind()
		return count
	}

	l.Rewind()

	return 0
}

func peekString(l *lexer.L, count int) string {
	s := ""

	for i := 0; i < count; i++ {
		if l.Peek() == lexer.EOFRune {
			return s
		}

		s += string(l.Next())
		defer l.Rewind()
	}

	return s
}
