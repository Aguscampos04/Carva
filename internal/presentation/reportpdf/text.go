package reportpdf

import (
	"fmt"
	"strings"
	"unicode"
)

// safeText keeps all visible information: glyphs absent from the embedded font
// are represented by their Unicode code point instead of silently disappearing.
func safeText(glyphs map[int]uint, text string) string {
	var output strings.Builder
	for _, r := range text {
		switch {
		case r == '\n':
			output.WriteRune(r)
		case r == '\t':
			output.WriteByte(' ')
		case unicode.IsControl(r):
			continue
		case glyphs[int(r)] == 0:
			fmt.Fprintf(&output, "[U+%04X]", r)
		default:
			output.WriteRune(r)
		}
	}
	return output.String()
}
