package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"unicode/utf8"

	"github.com/Aguscampos04/Carva/internal/domain/member"
)

// MemberService exposes registration and listing without coupling HTTP to SQLite.
type MemberService interface {
	Create(context.Context, string, string) (member.Member, error)
	List(context.Context, string) ([]member.Member, error)
}

func registerMemberRoutes(mux *http.ServeMux, service MemberService) {
	mux.HandleFunc("POST /api/projects/{projectID}/members", func(w http.ResponseWriter, r *http.Request) {
		name, err := decodeMemberName(r.Body)
		if err != nil {
			writeError(w, 400, "invalid JSON request")
			return
		}
		created, err := service.Create(r.Context(), r.PathValue("projectID"), name)
		if memberServiceError(w, err, "could not register member") {
			return
		}
		writeJSON(w, http.StatusCreated, created)
	})
	mux.HandleFunc("GET /api/projects/{projectID}/members", func(w http.ResponseWriter, r *http.Request) {
		items, err := service.List(r.Context(), r.PathValue("projectID"))
		if memberServiceError(w, err, "could not retrieve members") {
			return
		}
		writeJSON(w, http.StatusOK, items)
	})
}

func memberServiceError(w http.ResponseWriter, err error, fallback string) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, member.ErrInvalidName):
		writeError(w, 400, "member name must not be blank")
	case errors.Is(err, member.ErrProjectNotFound):
		writeError(w, 404, "project not found")
	default:
		writeError(w, 500, fallback)
	}
	return true
}

// decodeMemberName validates the entire object, exact keys, duplicates and EOF.
// Unknown fields are parsed but ignored; they cannot set server-controlled fields.
func decodeMemberName(body io.Reader) (string, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}
	invalid := errors.New("invalid member JSON object")
	if !utf8.Valid(data) {
		return "", invalid
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	token, err := decoder.Token()
	if err != nil || token != json.Delim('{') {
		return "", invalid
	}
	seen := map[string]bool{}
	var name string
	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			return "", invalid
		}
		key, ok := token.(string)
		if !ok || seen[key] {
			return "", invalid
		}
		seen[key] = true
		var value json.RawMessage
		if err := decoder.Decode(&value); err != nil {
			return "", invalid
		}
		if key == "name" {
			if err := json.Unmarshal(value, &name); err != nil {
				return "", invalid
			}
		}
	}
	if token, err := decoder.Token(); err != nil || token != json.Delim('}') {
		return "", invalid
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return "", invalid
	}
	return name, nil
}
