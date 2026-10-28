package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	app "github.com/Aguscampos04/Carva/internal/application/backlog"
	domain "github.com/Aguscampos04/Carva/internal/domain/backlog"
)

func WithAssignees(base http.Handler, service *app.AssigneeService) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", base)
	handle := func(w http.ResponseWriter, r *http.Request) {
		var value app.Assignee
		var err error
		pid, id := r.PathValue("projectID"), r.PathValue("itemID")
		if r.Method == http.MethodPut {
			var fields map[string]json.RawMessage
			if decodeSingleJSON(r, &fields) != nil || len(fields) != 1 || fields["member_id"] == nil {
				writeError(w, 400, "member_id is required; use null to unassign")
				return
			}
			var memberID *string
			if json.Unmarshal(fields["member_id"], &memberID) != nil {
				writeError(w, 400, "invalid member_id")
				return
			}
			value, err = service.Set(r.Context(), pid, id, memberID)
		} else {
			value, err = service.Get(r.Context(), pid, id)
		}
		switch {
		case errors.Is(err, domain.ErrNotFound):
			writeError(w, 404, "story not found")
		case errors.Is(err, app.ErrInvalidAssignee):
			writeError(w, 400, "el responsable debe ser un miembro del mismo proyecto")
		case err != nil:
			writeError(w, 500, "could not access story assignee")
		default:
			writeJSON(w, 200, value)
		}
	}
	mux.HandleFunc("GET /api/projects/{projectID}/backlog/{itemID}/assignee", handle)
	mux.HandleFunc("PUT /api/projects/{projectID}/backlog/{itemID}/assignee", handle)
	return mux
}
