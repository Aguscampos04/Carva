ALTER TABLE sprint_backlog_items
    ADD COLUMN completion_verified INTEGER NOT NULL DEFAULT 0
    CHECK (completion_verified IN (0, 1));
