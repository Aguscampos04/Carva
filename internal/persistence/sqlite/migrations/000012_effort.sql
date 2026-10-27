-- US-16: preserve legacy null estimates and existing planning/results.
ALTER TABLE backlog_items ADD COLUMN estimated_hours_hundredths INTEGER NULL
    CHECK (estimated_hours_hundredths IS NULL OR (typeof(estimated_hours_hundredths) = 'integer' AND estimated_hours_hundredths > 0));
ALTER TABLE sprint_backlog_items ADD COLUMN estimated_hours_hundredths INTEGER NULL
    CHECK (estimated_hours_hundredths IS NULL OR (typeof(estimated_hours_hundredths) = 'integer' AND estimated_hours_hundredths > 0));
CREATE UNIQUE INDEX backlog_project_identity ON backlog_items(project_id,id);
CREATE UNIQUE INDEX member_project_identity ON members(project_id,id);
CREATE UNIQUE INDEX sprint_project_identity ON sprints(project_id,id);
CREATE TABLE effort_entries (
    id TEXT PRIMARY KEY NOT NULL,
    project_id TEXT NOT NULL REFERENCES projects(id),
    backlog_item_id TEXT NOT NULL,
    member_id TEXT NOT NULL,
    work_date TEXT NOT NULL,
    activity TEXT NOT NULL CHECK(length(CAST(activity AS BLOB)) > 0),
    hours_hundredths INTEGER NOT NULL CHECK(typeof(hours_hundredths) = 'integer' AND hours_hundredths > 0 AND hours_hundredths <= 2400),
    sprint_id INTEGER NULL,
    deleted_at TEXT NULL,
    FOREIGN KEY(project_id,backlog_item_id) REFERENCES backlog_items(project_id,id),
    FOREIGN KEY(project_id,member_id) REFERENCES members(project_id,id),
    FOREIGN KEY(project_id,sprint_id) REFERENCES sprints(project_id,id),
    FOREIGN KEY(sprint_id,backlog_item_id) REFERENCES sprint_backlog_items(sprint_id,backlog_item_id)
);
CREATE INDEX effort_story_active ON effort_entries(project_id,backlog_item_id,deleted_at);
CREATE INDEX effort_daily_active ON effort_entries(project_id,member_id,work_date,deleted_at);
