package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"text/template"
	"time"
)

const versionInfoTmpl = `
{{.Program}}, version {{.Version}} (revision: {{.Commit}})
  build date:       {{.Date}}
  go version:       {{.GoVersion}}
`

// BuildInfo contains information about build
// like version tag, commit and build date.
type BuildInfo struct {
	version        string
	date           time.Time
	commit         string
	runtimeVersion func() string
}

// NewBuildInfo creates a new BuildInfo. The date has to be in RFC3329 format
// and the commit hash has to be at leas 8 characters long.
func NewBuildInfo(version, date, commit string) (*BuildInfo, error) {
	if version == "" {
		return nil, errors.New("version cannot be empty")
	}

	d, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date %s: %w", date, err)
	}

	if len(commit) < 8 {
		return nil, errors.New("commit hash has to be at least 8 characters long")
	}

	return &BuildInfo{
		version:        version,
		date:           d,
		commit:         commit,
		runtimeVersion: runtime.Version,
	}, nil
}

// Version returns the version information as string.
func (b BuildInfo) Version(program string) string {
	t := template.Must(template.New("version").Parse(versionInfoTmpl))

	data := struct {
		Version   string
		Date      string
		Commit    string
		GoVersion string
		Program   string
	}{
		Version:   b.version,
		Date:      b.date.In(time.UTC).String(),
		Commit:    b.commit[:8],
		GoVersion: b.runtimeVersion(),
		Program:   program,
	}

	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "version", data); err != nil {
		panic(err)
	}

	return strings.TrimSpace(buf.String())
}
