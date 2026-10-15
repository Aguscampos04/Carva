CREATE TABLE projects (
    id TEXT PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('Planificado', 'En progreso', 'Finalizado'))
);
