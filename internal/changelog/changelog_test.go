package changelog

import (
	"bytes"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrite(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	err = c.AddMessage("00000001", `fix: a fix

this is the body`)
	require.NoError(t, err)
	err = c.AddMessage("00000002", "fix: fixed this")
	require.NoError(t, err)
	err = c.AddMessage("00000003", `fix(router): another fix

Closes: #1`)
	require.NoError(t, err)
	err = c.AddMessage("00000004", "feat: a feature")
	require.NoError(t, err)
	err = c.AddMessage("00000005", `chore(common): changed this

Closes: #123`)
	require.NoError(t, err)
	err = c.AddMessage("00000006", "chore!: breaking change")
	require.NoError(t, err)
	err = c.AddMessage("00000007", "chore(router)!: breaking change again")
	require.NoError(t, err)
	err = c.AddMessage("00000008", `feat(router): add a breaking change

this is the body of the breaking change

BREAKING CHANGE: breaks all
Closes: #12345`)
	require.NoError(t, err)

	expected := `## title


### Breaking Changes

* **common**
  * **00000006**:
    breaking change
* **router**
  * **00000008**:
    add a breaking change (#12345)
    > this is the body of the breaking change
    > breaks all
  * **00000007**:
    breaking change again


### Bug Fixes

* **common**: a fix (00000001)
  > this is the body
* **common**: fixed this (00000002)
* **router**: another fix (#1, 00000003)


### New Features

* **common**: a feature (00000004)



`

	b := bytes.NewBufferString("")
	c.Write("title", b)

	if runtime.GOOS == "windows" {
		expected = strings.ReplaceAll(expected, "\n", "\r\n")
	}

	assert.Equal(t, expected, b.String())
}

func TestReleaseType(t *testing.T) {
	var tt = []struct {
		name     string
		messages []string
		expected ReleaseType
	}{
		{
			"patch",
			[]string{"chore: chore"},
			Patch,
		},
		{
			"patch",
			[]string{"chore: chore", "fix: fix"},
			Patch,
		},
		{
			"minor",
			[]string{"chore: chore", "feat: feat", "fix: fix"},
			Minor,
		},
		{
			"major",
			[]string{"chore!: chore", "feat: feat", "fix: fix"},
			Major,
		},
	}

	for i := range tt {
		tc := tt[i]

		t.Run(tc.name, func(t *testing.T) {
			cw, err := New()
			require.NoError(t, err)

			for _, m := range tc.messages {
				err := cw.AddMessage("aff5b0e55c1ede1c33425568f842e908f97eff89", m)
				require.NoError(t, err)
			}
			assert.Equal(t, tc.expected, cw.ReleaseType())
		})
	}
}
