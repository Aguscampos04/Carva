package httpapi

import (
	"bytes"
	"encoding/json"
	"io"

	domain "github.com/Aguscampos04/Carva/internal/domain/defect"
)

func decodeDefectPatch(body io.Reader) (domain.Patch, error) {
	fields, err := decodeProjectObject(body, "status", "detection_sprint_id", "resolution_sprint_id")
	if err != nil {
		return domain.Patch{}, err
	}
	var change domain.Patch
	if raw, ok := fields["status"]; ok {
		if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
			return change, domain.ErrInvalidStatus
		}
		var status domain.Status
		if err := json.Unmarshal(raw, &status); err != nil {
			return change, err
		}
		change.Status = &status
	}
	if raw, ok := fields["detection_sprint_id"]; ok {
		change.HasDetectionSprint = true
		if err := json.Unmarshal(raw, &change.DetectionSprintID); err != nil {
			return change, err
		}
	}
	if raw, ok := fields["resolution_sprint_id"]; ok {
		change.HasResolutionSprint = true
		if err := json.Unmarshal(raw, &change.ResolutionSprintID); err != nil {
			return change, err
		}
	}
	return change, change.Validate()
}
