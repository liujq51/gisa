package common

import (
	"strings"

	"github.com/beego/i18n"
)

func T(lang, format string, args ...interface{}) string {
	format = strings.ToLower(format)
	return i18n.Tr(lang, format, args)
}
