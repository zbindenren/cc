package cmd

import (
	"path/filepath"

	"github.com/postfinance/flash"
	"github.com/zbindenren/cc/config"
)

func (c Command) runWriteConfig(l *flash.Logger) error {
	l.Debugw("writing default configuration file", "path", filepath.Join(".", config.FileName))

	return config.Write(".", config.Default)
}
