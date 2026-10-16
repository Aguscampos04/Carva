CREATE TABLE members (
    id TEXT PRIMARY KEY NOT NULL,
    project_id TEXT NOT NULL REFERENCES projects(id),
    name TEXT NOT NULL CHECK (length(name) > 0)
);

CREATE INDEX members_project_id ON members (project_id);
