package httpapi

import (
	"errors"
	"fmt"
	app "github.com/Aguscampos04/Carva/internal/application/reports"
	"github.com/Aguscampos04/Carva/internal/domain/project"
	"github.com/Aguscampos04/Carva/internal/domain/sprint"
	"github.com/Aguscampos04/Carva/internal/presentation/reportpdf"
	"net/http"
	"strconv"
	"strings"
)

func WithReports(base http.Handler, service *app.Service) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", base)
	handler := func(w http.ResponseWriter, r *http.Request) {
		pid := r.PathValue("projectID")
		raw := r.PathValue("reportID")
		var result app.Report
		var err error
		pdf := strings.HasSuffix(raw, ".pdf")
		id := strings.TrimSuffix(raw, ".pdf")
		isSprint := strings.Contains(r.Pattern, "/sprints/")
		if !isSprint && id != "project" {
			writeError(w, 404, "report not found")
			return
		}
		filename := "carva-proyecto.pdf"
		if !isSprint {
			result, err = service.Project(r.Context(), pid)
		} else {
			var sid int64
			sid, err = strconv.ParseInt(id, 10, 64)
			if err != nil || sid < 1 {
				writeError(w, 400, app.ErrInvalidSprintID.Error())
				return
			}
			filename = fmt.Sprintf("carva-sprint-%d.pdf", sid)
			result, err = service.Sprint(r.Context(), pid, sid)
		}
		if err != nil {
			code := 500
			if errors.Is(err, project.ErrNotFound) || errors.Is(err, sprint.ErrNotFound) {
				code = 404
			}
			writeError(w, code, "could not generate requested report")
			return
		}
		if !pdf {
			writeJSON(w, 200, struct {
				app.Report
				Sections []reportpdf.Section `json:"sections"`
			}{result, reportpdf.Sections(result)})
			return
		}
		data, err := reportpdf.Render(result)
		if err != nil {
			writeError(w, 500, "could not generate PDF")
			return
		}
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(200)
		_, _ = w.Write(data)
	}
	// A single wildcard accepts both the numeric ID and its .pdf suffix.
	mux.HandleFunc("GET /api/projects/{projectID}/reports/{reportID}", handler)
	mux.HandleFunc("GET /api/projects/{projectID}/reports/sprints/{reportID}", handler)
	return mux
}
