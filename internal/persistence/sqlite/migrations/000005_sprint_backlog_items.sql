CREATE TABLE sprint_backlog_items (
    id INTEGER PRIMARY KEY,
    sprint_id INTEGER NOT NULL REFERENCES sprints(id),
    backlog_item_id TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    priority TEXT NOT NULL,
    status TEXT NOT NULL,
    provisional_story_points INTEGER CHECK (provisional_story_points IS NULL OR provisional_story_points >= 0),
    story_points INTEGER CHECK (story_points IS NULL OR story_points >= 0),
    acceptance_criteria TEXT NOT NULL,
    UNIQUE (sprint_id, backlog_item_id)
);

CREATE INDEX sprint_backlog_items_sprint_id_idx ON sprint_backlog_items(sprint_id);
