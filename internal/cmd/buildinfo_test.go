package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func TestBuildInfo(t *testing.T) {
	var tt = []struct {
		name         string
		version      string
		commit       string
		date         string
		expected     string
		expectsError bool
	}{
		{
			"ok",
			"1.0.0",
			"4a8804f18ea7560fe45fecaa052605b1a8a66fe8",
			"2021-01-05T08:08:54+01:00",
			`prog, version 1.0.0 (revision: 4a8804f1)
  build date:       2021-01-05 08:08:54 +0100 CET
  go version:       go1.15.6`,
			false,
		},
		{
			"invalid - commit to short",
			"1.0.0",
			"123",
			"2021-01-05T08:08:54+01:00",
			"",
			true,
		},
		{
			"invalid - empty version",
			"",
			"4a8804f18ea7560fe45fecaa052605b1a8a66fe8",
			"2021-01-05T08:08:54+01:00",
			"",
			true,
		},
		{
			"invalid - wrong date format",
			"1.0.0",
			"4a8804f18ea7560fe45fecaa052605b1a8a66fe8",
			"2021-01-05T08:08:54",
			"",
			true,
		},
	}

	for i := range tt {
		tc := tt[i]
		t.Run(tc.name, func(t *testing.T) {
			b, err := NewBuildInfo(tc.version, tc.date, tc.commit)
			if tc.expectsError {
				require.Error(t, err)
				return
			}
			b.runtimeVersion = runtimeVersionMock
			require.NoError(t, err)

			assert.Equal(t, tc.expected, b.Version("prog"))
		})
	}
}

func runtimeVersionMock() string {
	return "go1.15.6"
}
