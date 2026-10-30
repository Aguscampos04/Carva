package reportpdf

import (
	"encoding/json"
	"fmt"
	app "github.com/Aguscampos04/Carva/internal/application/reports"
	"strings"
)

type Section struct {
	Title string   `json:"title"`
	Lines []string `json:"lines"`
}

func numeric(value any) string {
	data, err := json.Marshal(value)
	if err != nil || string(data) == "null" {
		return "No disponible"
	}
	return string(data)
}
func optional(value string) string {
	if value == "" {
		return "Sin registrar"
	}
	return value
}
func sprintID(id *int64) string {
	if id == nil {
		return "Sin Sprint"
	}
	return fmt.Sprint(*id)
}
func metricLines(value app.Summary) []string {
	return []string{
		"Story Points planificados: " + numeric(value.Points.PlannedStoryPoints),
		"Story Points completados: " + numeric(value.Points.CompletedStoryPoints),
		"Velocidad del equipo: " + numeric(value.Points.TeamVelocity),
		"Historias completadas (%): " + numeric(value.Points.CompletionPercentage),
		"Horas estimadas: " + numeric(value.Hours.EstimatedHours) + " h",
		"Horas reales: " + numeric(value.Hours.ActualHours) + " h",
		"Desviación (reales menos estimadas): " + numeric(value.Hours.Deviation) + " h",
		"Defectos detectados: " + numeric(value.Defects.DetectedDefects),
		"Defectos resueltos: " + numeric(value.Defects.ResolvedDefects),
	}
}

// Sections is a textual projection of the same report returned as JSON.
// It never calculates totals or reconstructs snapshots from current Backlog.
func Sections(report app.Report) []Section {
	result := []Section{{"Información general", []string{
		"Proyecto: " + report.Project.Name, "Identificador: " + report.Project.ID,
		"Estado: " + string(report.Project.Status), "Inicio: " + optional(string(report.Project.StartDate)),
		"Finalización: " + optional(string(report.Project.EndDate)),
		"Alcance: " + map[string]string{"project": "Proyecto completo", "sprint": "Sprint individual"}[report.Scope],
		"Datos consultados al generar el reporte; no certifica aceptación ni cierre.",
	}}, {"Métricas del alcance", metricLines(report.Metrics)}}
	members := Section{Title: "Miembros del proyecto"}
	for _, m := range report.Members {
		members.Lines = append(members.Lines, m.Name+" — "+m.ID)
	}
	result = append(result, members)
	if report.Scope == "project" {
		current := Section{Title: "Historias del Product Backlog actual"}
		for _, s := range report.Stories {
			current.Lines = append(current.Lines, s.Title+" — "+s.ID, "Estado: "+string(s.Status)+"; prioridad: "+string(s.Priority), s.Description,
				"Story Points definitivos: "+numeric(s.StoryPoints)+"; provisionales: "+numeric(s.ProvisionalStoryPoints),
				"Horas estimadas vigentes: "+numeric(s.EstimatedHours)+" h; reales: "+numeric(s.ActualHours)+" h",
				"Criterios: "+strings.Join(s.AcceptanceCriteria, "; "))
		}
		result = append(result, current)
	}
	if len(report.Sprints) == 0 {
		result = append(result, Section{"Sprints", nil})
	}
	for _, value := range report.Sprints {
		state := "Abierto"
		if value.Sprint.Closed {
			state = "Cerrado"
		}
		result = append(result, Section{fmt.Sprintf("Sprint %d — %s", value.Sprint.ID, value.Sprint.Goal), append([]string{"Estado: " + state}, metricLines(value.Metrics)...)})
		planned := Section{Title: fmt.Sprintf("Historias planificadas y resultados — Sprint %d", value.Sprint.ID)}
		for _, s := range value.Stories {
			outcome := "Incompleta"
			if s.AcceptanceCriteriaVerified {
				outcome = "Completada"
			}
			planned.Lines = append(planned.Lines, s.Title+" — "+s.BacklogItemID, s.Description,
				"Estado planificado: "+s.PlannedStatus+"; estado registrado: "+s.Status+"; resultado: "+outcome,
				"Prioridad: "+s.Priority+"; Story Points definitivos: "+numeric(s.StoryPoints)+"; provisionales: "+numeric(s.ProvisionalStoryPoints),
				"Horas estimadas vigentes (US-16): "+numeric(s.EstimatedHours)+" h",
				"Criterios: "+strings.Join(s.AcceptanceCriteria, "; "))
		}
		result = append(result, planned)
	}
	entries := Section{Title: "Esfuerzo registrado activo"}
	for _, e := range report.Effort {
		entries.Lines = append(entries.Lines, fmt.Sprintf("%s — miembro %s; historia %s; Sprint %s; %s h", e.WorkDate, e.MemberID, e.StoryID, sprintID(e.SprintID), numeric(e.Hours)), e.Activity)
	}
	result = append(result, entries)
	defects := Section{Title: "Defectos"}
	for _, d := range report.Defects {
		defects.Lines = append(defects.Lines, d.ID+" — "+d.Description,
			"Historia: "+d.BacklogItemID+"; severidad: "+string(d.Severity)+"; estado: "+string(d.Status),
			"Sprint de detección: "+sprintID(d.DetectionSprintID)+"; resolución: "+sprintID(d.ResolutionSprintID))
	}
	result = append(result, defects)
	return result
}
