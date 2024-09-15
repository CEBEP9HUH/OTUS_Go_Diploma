package util

import (
	"strconv"
	"strings"
)

func GetFloat(data string) (float32, error) {
	if data == "-" {
		return 0, nil
	}
	data = strings.TrimSpace(data)
	v, err := strconv.ParseFloat(strings.Replace(data, ",", ".", 1), 32)
	return float32(v), err
}

func GetTrimmedFloat(data, prefix, suffix string) (float32, error) {
	value := strings.TrimPrefix(data, prefix)
	value = strings.TrimSuffix(value, suffix)
	return GetFloat(value)
}
