CREATE TABLE projects (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL CHECK (length(trim(name)) > 0),
    status TEXT NOT NULL CHECK (status IN ('Planificado', 'En progreso', 'Finalizado'))
);

CREATE TABLE sprints (
    id INTEGER PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id),
    goal TEXT NOT NULL CHECK (length(trim(goal)) > 0)
);

CREATE INDEX sprints_project_id_idx ON sprints(project_id);
