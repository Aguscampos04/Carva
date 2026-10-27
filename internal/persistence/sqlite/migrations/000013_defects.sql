CREATE TABLE defects (
    id TEXT PRIMARY KEY NOT NULL,
    project_id TEXT NOT NULL REFERENCES projects(id),
    backlog_item_id TEXT NOT NULL,
    description TEXT NOT NULL CHECK(length(trim(description)) > 0),
    severity TEXT NOT NULL CHECK(severity IN ('Baja', 'Media', 'Alta', 'Crítica')),
    status TEXT NOT NULL CHECK(status IN ('Abierto', 'En progreso', 'Resuelto', 'Cerrado')),
    detection_sprint_id INTEGER NULL,
    resolution_sprint_id INTEGER NULL,
    FOREIGN KEY(project_id, backlog_item_id) REFERENCES backlog_items(project_id, id),
    FOREIGN KEY(project_id, detection_sprint_id) REFERENCES sprints(project_id, id),
    FOREIGN KEY(project_id, resolution_sprint_id) REFERENCES sprints(project_id, id)
);

CREATE INDEX defects_story_idx ON defects(project_id, backlog_item_id);
CREATE INDEX defects_detection_sprint_idx ON defects(project_id, detection_sprint_id);
CREATE INDEX defects_resolution_sprint_idx ON defects(project_id, resolution_sprint_id);
