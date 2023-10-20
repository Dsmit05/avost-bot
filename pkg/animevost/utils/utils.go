package utils

import (
	"regexp"
	"strings"
)

var reg = regexp.MustCompile(`\w+`)

func GetTitleNameForURL(s string) string {
	start := strings.Index(s, "/")
	end := strings.Index(s, "[")

	if start == -1 || end == -1 {
		return ""
	}

	if start >= end {
		return "" // ToDo: Invalid index values, must be low <= high
	}

	elems := reg.FindAllString(s[start+2:end], -1)

	var buf strings.Builder

	sep := "-"
	n := len(sep) * (len(elems) - 1)

	for i := 0; i < len(elems); i++ {
		n += len(elems[i])
	}

	buf.Grow(n)
	buf.WriteString(strings.ToLower(elems[0]))

	for _, s := range elems[1:] {
		buf.WriteString(sep)
		buf.WriteString(strings.ToLower(s))
	}

	return buf.String()
}
