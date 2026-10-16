// Package member defines project member identities and name validation (US-03).
package member

import (
	"errors"
	"strings"
)

var (
	ErrInvalidName     = errors.New("member name must not be blank")
	ErrInvalidIdentity = errors.New("member identity and project must not be blank")
	ErrProjectNotFound = errors.New("project not found")
)

// Member is a registration within one project, independent of other registrations.
type Member struct {
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
	Name      string `json:"name"`
}

// New preserves the name's interior and trims only surrounding whitespace.
func New(id, projectID, name string) (Member, error) {
	name, err := NormalizeName(name)
	if err != nil {
		return Member{}, err
	}
	if strings.TrimSpace(id) == "" || strings.TrimSpace(projectID) == "" {
		return Member{}, ErrInvalidIdentity
	}
	return Member{ID: id, ProjectID: projectID, Name: name}, nil
}

// NormalizeName is the single definition of D-06, also used before project lookup.
func NormalizeName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", ErrInvalidName
	}
	return name, nil
}
