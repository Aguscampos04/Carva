CREATE TABLE story_assignees (
    backlog_item_id TEXT PRIMARY KEY NOT NULL REFERENCES backlog_items(id),
    member_id TEXT NOT NULL REFERENCES members(id)
);
