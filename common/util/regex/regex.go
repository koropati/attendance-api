package regex

import (
	"regexp"
	"strings"
)

const (
	NAME = "^[a-zA-Z\\s]{2,40}$"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(camelCase string) string {
	snake := matchFirstCap.ReplaceAllString(camelCase, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
