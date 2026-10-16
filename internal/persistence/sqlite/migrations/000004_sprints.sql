CREATE TABLE sprints (
    id INTEGER PRIMARY KEY,
    project_id TEXT NOT NULL REFERENCES projects(id),
    goal TEXT NOT NULL CHECK (length(trim(goal)) > 0)
);

CREATE INDEX sprints_project_id_idx ON sprints(project_id);