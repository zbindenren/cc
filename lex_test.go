package cc

import (
	"testing"

	"github.com/bbuck/go-lexer"
	"github.com/tj/assert"
)

func TestTakeUntilFirstFooter(t *testing.T) {
	tt := []struct {
		data            string
		expectedFound   bool
		expectedCurrent string
	}{
		{
			"some data\nmore data\n\nfooter: some value\n",
			true,
			"some data\nmore data\n\n",
		},
		{
			"some data\nmore data\nfooter: some value\n",
			true,
			"some data\nmore data\n",
		},
		{
			"some data\nmore data\nBREAKING-CHANGE: some value\n",
			true,
			"some data\nmore data\n",
		},
		{
			"some data\nbreaks\nBREAKING CHANGE: the breaking change\n",
			true,
			"some data\nbreaks\n",
		},
		{
			"some data\nmore data\n\n\nfooter: some value\n",
			true,
			"some data\nmore data\n\n\n",
		},
		{
			"some data\nmore data footer: some value\n",
			false,
			"some data\nmore data footer: some value\n",
		},
		{
			"\n\n\nfooter: some value",
			true,
			"\n\n\n",
		},
	}

	for i := range tt {
		tc := tt[i]

		t.Run(tc.data, func(t *testing.T) {
			l := lexer.New(tc.data, nil)
			found := takeUntilFirstFooterToken(l)
			assert.Equal(t, tc.expectedFound, found)
			assert.Equal(t, tc.expectedCurrent, l.Current())
		})
	}
}

func TestTakeFooterToken(t *testing.T) {
	tt := []struct {
		data          string
		expectedCount int
		current       string
	}{
		{
			"footer: some value",
			6,
			"footer",
		},
		{
			"footer #some value",
			6,
			"footer",
		},
		{
			"Boot: some value",
			4,
			"Boot",
		},
		{
			"BREAKING-CHANGE: some value",
			15,
			"BREAKING-CHANGE",
		},
		{
			"BREAKING CHANGE: some value",
			15,
			"BREAKING CHANGE",
		},
		{
			"BREAKING_CHANGE: some value",
			0,
			"BREAKING_CHANGE",
		},
		{
			"BREAKING CHANGE some value",
			0,
			"BREAKING CHANGE",
		},
		{
			"no footer: some value",
			0,
			"no",
		},
		{
			"no footer #some value",
			0,
			"no",
		},
		{
			"\nfooter: some value",
			0,
			"\n",
		},
	}

	for i := range tt {
		tc := tt[i]

		t.Run(tc.data, func(t *testing.T) {
			l := lexer.New(tc.data, nil)
			count := takeFooterToken(l)
			assert.Equal(t, tc.expectedCount, count)
		})
	}
}
