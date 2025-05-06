package utils_test

import (
	"testing"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestSanitizeString(t *testing.T) {
	raw := `<script>alert("xss")</script><p>Hello <b>World</b></p>`
	expected := `Hello World` // StrictPolicy удаляет все теги

	clean := utils.SanitizeString(raw)
	assert.Equal(t, expected, clean)
}

func TestSanitizeRichText(t *testing.T) {
	raw := `<script>alert("xss")</script><p>Hello <b>World</b></p><a href="http://example.com">Link</a>`
	expected := `<p>Hello <b>World</b></p><a href="http://example.com" rel="nofollow">Link</a>`

	clean := utils.SanitizeRichText(raw)
	assert.Equal(t, expected, clean)
}

func TestSanitizeEmptyString(t *testing.T) {
	assert.Equal(t, "", utils.SanitizeString(""))
	assert.Equal(t, "", utils.SanitizeRichText(""))
}
