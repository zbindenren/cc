package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/postfinance/flash"
	"github.com/stretchr/testify/require"
	"github.com/zbindenren/cc/internal/git"
	"gotest.tools/assert"
)

const (
	logOutput = false
	windowsOS = "windows"
)

func TestInit(t *testing.T) {
	c, changlogPath, cleanup := setup(t, "untagged")
	defer cleanup()

	err := c.Run()
	require.NoError(t, err)

	b, err := ioutil.ReadFile(changlogPath) // nolint: gosec
	require.NoError(t, err)

	expectedFmt := `## 0.1.0 (%s)


### Bug Fixes

* **common**: fix an error (0fec975c9da5c5ce62f63c9d7bc0009255451006)
  > this is the body of the message.
  > can be multiline.


### New Features

* **common**: initial working version (c49e021712062196bff430c0acff8312dc343b74)
* **router**: add new router flag (596ae7e44b5bfb2792d237b29159c5cc51a10a25)



`
	expected := fmt.Sprintf(expectedFmt, time.Now().Format(dateFormat))

	// '\n' vs '\r\n'
	if runtime.GOOS == windowsOS {
		expected = strings.ReplaceAll(expected, "\n", "\r\n")
	}

	assert.Equal(t, expected, string(b))
}

func TestRelease(t *testing.T) {
	c, changlogPath, cleanup := setup(t, "tagged")
	defer cleanup()

	err := c.Run()
	require.NoError(t, err)

	b, err := ioutil.ReadFile(changlogPath) // nolint: gosec
	require.NoError(t, err)

	expectedFmt := `## 0.2.0 (%s)


### Bug Fixes

* **common**: fix error handling (aa5b93a9ee73be410eeab92d4276b257d15ecf6b)


### New Features

* **common**: add new feature (14f3b06858668ef50ccbcccf8266f495b434f71c)



## 0.1.0 (2020-12-30)


### Bug Fixes

* **common**: fix an error (0fec975c9da5c5ce62f63c9d7bc0009255451006)
  > this is the body of the message.
  > can be multiline.


### New Features

* **common**: initial working version (c49e021712062196bff430c0acff8312dc343b74)
* **router**: add new router flag (596ae7e44b5bfb2792d237b29159c5cc51a10a25)



`
	expected := fmt.Sprintf(expectedFmt, time.Now().Format(dateFormat))

	// '\n' vs '\r\n'
	if runtime.GOOS == windowsOS {
		expected = strings.ReplaceAll(expected, "\n", "\r\n")
	}

	assert.Equal(t, expected, string(b))
}

func setup(t *testing.T, repoName string) (c Command, changelogPath string, cleanup func()) {
	tmp, err := ioutil.TempDir("", repoName)
	require.NoError(t, err)

	old, err := os.Getwd()
	require.NoError(t, err)

	repoDir := filepath.Join(tmp, repoName)

	g, err := git.New(flash.New())
	require.NoError(t, err)
	_, err = g.Run("clone", pathToBundle(repoName), repoDir)
	require.NoError(t, err)

	err = os.Chdir(repoDir)
	require.NoError(t, err)

	c = Command{
		noop:       true,
		debug:      newBoolPtr(logOutput),
		initConfig: newBoolPtr(false),
		history:    newBoolPtr(false),
		noPrompt:   newBoolPtr(true),
		toStdOut:   newBoolPtr(false),
		file:       newStrPtr(dfltChangelogFile),
		version:    newBoolPtr(false),
		num:        newIntPtr(0),
		sinceTag:   newStrPtr(""),
	}

	cleanup = func() {
		os.RemoveAll(tmp)

		if err := os.Chdir(old); err != nil {
			panic(err)
		}
	}

	return c, filepath.Join(tmp, repoName, dfltChangelogFile), cleanup
}

func newBoolPtr(b bool) *bool {
	return &b
}

func newIntPtr(i int) *int {
	return &i
}

func newStrPtr(s string) *string {
	return &s
}

func pathToBundle(name string) string {
	_, filename, _, _ := runtime.Caller(0) // nolint: dogsled
	dir := filepath.Dir(filename)

	return filepath.Join(dir, "test", name+".bundle")
}
