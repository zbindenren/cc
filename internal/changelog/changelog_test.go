package changelog

import (
	"bytes"
	"io/ioutil"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

const windowsOS = "windows"

// nolint: funlen
func TestWrite(t *testing.T) {
	t.Run("gitlab", func(t *testing.T) {
		c, err := New()
		require.NoError(t, err)

		messages := messagesFrom(t, "test-messages.yml")

		for _, m := range messages {
			err := c.AddMessage(m.Commit, m.Message)
			require.NoError(t, err)
		}

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

		if runtime.GOOS == windowsOS {
			expected = strings.ReplaceAll(expected, "\n", "\r\n")
		}

		assert.Equal(t, expected, b.String())
	})

	t.Run("github", func(t *testing.T) {
		c, err := New()
		require.NoError(t, err)
		c.cfg.GithubProjectPath = "zbindenren/cc"

		messages := messagesFrom(t, "test-messages.yml")

		for _, m := range messages {
			err := c.AddMessage(m.Commit, m.Message)
			require.NoError(t, err)
		}

		expected := `## title


### Breaking Changes

* **common**
  * **[00000006](https://github.com/zbindenren/cc/commit/00000006)**:
    breaking change
* **router**
  * **[00000008](https://github.com/zbindenren/cc/commit/00000008)**:
    add a breaking change ([#12345](https://github.com/zbindenren/cc/issues/12345))
    > this is the body of the breaking change
    > breaks all
  * **[00000007](https://github.com/zbindenren/cc/commit/00000007)**:
    breaking change again


### Bug Fixes

* **common**: a fix ([00000001](https://github.com/zbindenren/cc/commit/00000001))
  > this is the body
* **common**: fixed this ([00000002](https://github.com/zbindenren/cc/commit/00000002))
* **router**: another fix ([#1](https://github.com/zbindenren/cc/issues/1), [00000003](https://github.com/zbindenren/cc/commit/00000003))


### New Features

* **common**: a feature ([00000004](https://github.com/zbindenren/cc/commit/00000004))



`

		b := bytes.NewBufferString("")
		c.Write("title", b)

		if runtime.GOOS == windowsOS {
			expected = strings.ReplaceAll(expected, "\n", "\r\n")
		}

		assert.Equal(t, expected, b.String())
	})
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

type message struct {
	Commit  string
	Message string
}

func messagesFrom(t *testing.T, filePath string) []message {
	m := []message{}

	d, err := ioutil.ReadFile(filePath) // nolint: gosec
	require.NoError(t, err)

	err = yaml.Unmarshal(d, &m)
	require.NoError(t, err)

	return m
}
