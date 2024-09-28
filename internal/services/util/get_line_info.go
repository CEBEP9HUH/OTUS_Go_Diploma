package util

import (
	"errors"
	"strings"
)

var ErrInfoLineNotFound = errors.New("info line not found")

func GetLineInfo(data, prefix, lineSep string) (string, error) {
	lData := strings.ToLower(data)
	lines := strings.Split(lData, lineSep)
	for _, l := range lines {
		if strings.Contains(l, prefix) {
			pos := strings.Index(l, prefix) + len(prefix)
			return l[pos:], nil
		}
	}
	return "", ErrInfoLineNotFound
}
