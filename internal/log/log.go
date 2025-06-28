package log

import (
	"github.com/go-pkgz/lgr"
)

func New(level string) lgr.L {
	return lgr.New(lgr.Msec, lgr.LevelBraces)
}
