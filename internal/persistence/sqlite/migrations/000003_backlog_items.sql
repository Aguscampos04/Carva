CREATE TABLE backlog_items (
    id TEXT PRIMARY KEY NOT NULL,
    project_id TEXT NOT NULL REFERENCES projects(id),
    title TEXT NOT NULL CHECK (length(trim(title)) > 0),
    description TEXT NOT NULL CHECK (length(trim(description)) > 0),
    priority TEXT NOT NULL CHECK (priority IN ('Alta', 'Media', 'Baja')),
    status TEXT NOT NULL DEFAULT 'Pendiente' CHECK (status IN ('Pendiente', 'En progreso', 'Completada')),
    provisional_story_points INTEGER CHECK (provisional_story_points IS NULL OR provisional_story_points >= 0),
    story_points INTEGER CHECK (story_points IS NULL OR story_points >= 0),
    acceptance_criteria TEXT NOT NULL
);

CREATE INDEX backlog_items_project_id ON backlog_items (project_id);
