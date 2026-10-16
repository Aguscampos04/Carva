-- SQLite length(TEXT) stops at NUL. D-06 only rejects blank names, so
-- enforce nonempty storage by byte length without restricting characters.
CREATE TABLE members_name_storage (
    id TEXT PRIMARY KEY NOT NULL,
    project_id TEXT NOT NULL REFERENCES projects(id),
    name TEXT NOT NULL CHECK (length(CAST(name AS BLOB)) > 0)
);

INSERT INTO members_name_storage (id, project_id, name)
SELECT id, project_id, name FROM members;

DROP TABLE members;
ALTER TABLE members_name_storage RENAME TO members;
CREATE INDEX members_project_id ON members (project_id);
