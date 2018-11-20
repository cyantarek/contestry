package utilities

import "strings"

func Slugify(s string) string {
	s = strings.Replace(s, " ", "-", -1)
	s = strings.ToLower(s)

	if len(s) > 15 {
		s = s[:15]
	}
	return s
}
