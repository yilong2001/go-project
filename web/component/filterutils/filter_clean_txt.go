package filterutils

import (
	"strings"
)

func CleanSearchTxt(txt string) string {
	search := strings.Replace(txt, "\\", "", -1)
	search = strings.Replace(search, "&", "", -1)
	search = strings.Replace(search, "<", "", -1)
	search = strings.Replace(search, ">", "", -1)
	search = strings.Replace(search, "/", "", -1)
	search = strings.Replace(search, "*", "", -1)
	search = strings.Replace(search, ".", "", -1)
	search = strings.Replace(search, "%", "", -1)
	search = strings.Replace(search, "$", "", -1)
	search = strings.Replace(search, "#", "", -1)
	search = strings.Replace(search, "@", "", -1)
	search = strings.Replace(search, "+", "", -1)
	search = strings.Replace(search, "-", "", -1)
	search = strings.Replace(search, "!", "", -1)
	search = strings.Replace(search, "~", "", -1)
	search = strings.Replace(search, "|", "", -1)
	search = strings.Replace(search, "(", "", -1)
	search = strings.Replace(search, ")", "", -1)
	search = strings.Replace(search, "=", "", -1)
	search = strings.Replace(search, "_", "", -1)
	search = strings.Replace(search, "'", "", -1)
	search = strings.Replace(search, "\"", "", -1)
	search = strings.Replace(search, "?", "", -1)

	return search
}
