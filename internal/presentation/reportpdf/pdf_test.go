package reportpdf_test

import (
 "bytes"
 "strings"
 "testing"
 app "github.com/Aguscampos04/Carva/internal/application/reports"
 "github.com/Aguscampos04/Carva/internal/domain/backlog"
 "github.com/Aguscampos04/Carva/internal/domain/project"
 pdf "github.com/Aguscampos04/Carva/internal/presentation/reportpdf"
)
func TestPDFPreservesSectionsUnavailableAndAccentedText(t *testing.T){
 report:=app.Report{Scope:"project",Project:project.Project{ID:"p",Name:"Estimación y resolución",Status:project.Planned},Stories:[]backlog.Item{{ID:"a",Title:strings.Repeat("Historia extensa con acentos áéíóú y ñ. ",150),Description:"última descripción",AcceptanceCriteria:[]string{"Criterio original"}}}}
 sections:=pdf.Sections(report)
 var lines []string
 for _,s:=range sections{lines=append(lines,s.Title);lines=append(lines,s.Lines...)}
 text:=strings.Join(lines,"\n")
 for _,want:=range []string{"Estimación y resolución","No disponible","Criterio original","última descripción","Métricas","Defectos","Esfuerzo","Historias"}{if !strings.Contains(text,want){t.Fatalf("missing %q in report",want)}}
 data,err:=pdf.Render(report);if err!=nil{t.Fatal(err)}
 if !bytes.HasPrefix(data,[]byte("%PDF-"))||!bytes.Contains(data,[]byte("%%EOF")){t.Fatal("invalid PDF")}
 if bytes.Count(data,[]byte("/Type /Page\n"))<2 {t.Fatal("long text must paginate")}
}
