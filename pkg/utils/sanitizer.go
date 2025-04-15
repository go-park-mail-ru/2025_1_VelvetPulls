// pkg/utils/sanitizer.go
package utils

import (
	"github.com/microcosm-cc/bluemonday"
)

type Sanitizable interface {
	Sanitize()
}

var (
	StrictPolicy = bluemonday.StrictPolicy()

	UGCPolicy = bluemonday.UGCPolicy()
)

func SanitizeString(input string) string {
	return StrictPolicy.Sanitize(input)
}

func SanitizeRichText(input string) string {
	return UGCPolicy.Sanitize(input)
}
