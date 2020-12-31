package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestWriteLoad(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "writeloadtest")
	require.NoError(t, err)

	defer os.RemoveAll(tmpDir)

	err = Write(tmpDir, Default)
	require.NoError(t, err)

	c, err := Load(tmpDir)
	require.NoError(t, err)

	require.Equal(t, Default, *c)

	_, err = Load(".")
	require.Equal(t, ErrNotFound, err)
}

func TestIsHidden(t *testing.T) {
	assert.True(t, Default.IsHidden("chore"))
	assert.False(t, Default.IsHidden("fix"))
}

func TestTitle(t *testing.T) {
	title, ok := Default.Title("feat")
	assert.True(t, ok)
	assert.Equal(t, "New Features", title)

	_, ok = Default.Title("not-exist")
	assert.False(t, ok)
}
