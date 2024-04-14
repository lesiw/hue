package main

import (
	"fmt"
	"strings"
)

var colors = map[string]string{
	"reset":     "0",
	"bold":      "1",
	"black":     "30",
	"red":       "31",
	"green":     "32",
	"yellow":    "33",
	"blue":      "34",
	"magenta":   "35",
	"cyan":      "36",
	"white":     "37",
	"default":   "39",
	"blackbg":   "40",
	"redbg":     "41",
	"greenbg":   "42",
	"yellowbg":  "43",
	"bluebg":    "44",
	"magentabg": "45",
	"cyanbg":    "46",
	"whitebg":   "47",
}

func colorToAnsi(s string) (string, error) {
	if len(s) < 1 {
		return "", nil
	}
	if s[0] == '\\' {
		return s, nil // Assume this is already an escape sequence.
	}
	var b strings.Builder
	b.WriteString("\033[")
	for i, c := range strings.Split(s, ",") {
		if i > 0 {
			b.WriteRune(';')
		}
		if s, ok := colors[c]; ok {
			b.WriteString(s)
		} else {
			return "", fmt.Errorf("bad color: '%s'", c)
		}
	}
	b.WriteRune('m')
	return b.String(), nil
}
