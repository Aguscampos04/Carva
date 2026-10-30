package reportpdf

import (
	"github.com/signintech/gopdf/fontmaker/core"
	"golang.org/x/image/font/gofont/goregular"
	"testing"
)

func TestPDFTextPreservesAccentsAndEscapesUnsupportedGlyphs(t *testing.T) {
	var font core.TTFParser
	if err := font.ParseFontData(goregular.TTF); err != nil {
		t.Fatal(err)
	}
	actual := safeText(font.Chars(), "Resolución áéíóú ñ\r\n\t🚀")
	if actual != "Resolución áéíóú ñ\n [U+1F680]" {
		t.Fatalf("text silently lost: %q", actual)
	}
}
