// Package reportpdf renders reports using an embedded redistributable Go font.
package reportpdf

import (
	"fmt"
	app "github.com/Aguscampos04/Carva/internal/application/reports"
	"github.com/signintech/gopdf"
	"golang.org/x/image/font/gofont/goregular"
	"strings"
	"unicode"
)

func Render(report app.Report) ([]byte, error) {
	var pdf gopdf.GoPdf
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	if err := pdf.AddTTFFontData("Go", goregular.TTF); err != nil {
		return nil, err
	}
	page := 0
	newPage := func() error {
		pdf.AddPage()
		page++
		if err := pdf.SetFont("Go", "", 10); err != nil {
			return err
		}
		pdf.SetXY(40, 25)
		if err := pdf.Cell(nil, fmt.Sprintf("Carva | Reporte | Página %d", page)); err != nil {
			return err
		}
		pdf.SetXY(40, 55)
		return nil
	}
	if err := newPage(); err != nil {
		return nil, err
	}
	write := func(text string, size int) error {
		if err := pdf.SetFont("Go", "", size); err != nil {
			return err
		}
		// Preserve newlines; normalize control characters unsuitable for PDF text.
		text = strings.Map(func(r rune) rune {
			if r == '\t' {
				return ' '
			}
			if r != '\n' && unicode.IsControl(r) {
				return -1
			}
			return r
		}, text)
		for _, paragraph := range strings.Split(text, "\n") {
			lines, err := pdf.SplitTextWithWordWrap(paragraph, 510)
			if err != nil {
				return err
			}
			if len(lines) == 0 {
				lines = []string{" "}
			}
			for _, line := range lines {
				if pdf.GetY()+20 > 795 {
					if err := newPage(); err != nil {
						return err
					}
					if err := pdf.SetFont("Go", "", size); err != nil {
						return err
					}
				}
				y := pdf.GetY()
				pdf.SetX(40)
				if err := pdf.Cell(nil, line); err != nil {
					return err
				}
				pdf.SetY(y + float64(size) + 5)
			}
		}
		return nil
	}
	for _, section := range Sections(report) {
		if pdf.GetY() > 745 {
			if err := newPage(); err != nil {
				return nil, err
			}
		}
		if err := write(section.Title, 14); err != nil {
			return nil, err
		}
		lines := section.Lines
		if len(lines) == 0 {
			lines = []string{"Sin registros."}
		}
		for _, line := range lines {
			if err := write(line, 10); err != nil {
				return nil, err
			}
		}
		pdf.SetY(pdf.GetY() + 10)
	}
	return pdf.GetBytesPdfReturnErr()
}
