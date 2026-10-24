CREATE TABLE estimation_rounds (
    id TEXT PRIMARY KEY NOT NULL,
    project_id TEXT NOT NULL REFERENCES projects(id),
    story_id TEXT NOT NULL REFERENCES backlog_items(id),
    round_number INTEGER NOT NULL CHECK (round_number > 0),
    status TEXT NOT NULL CHECK (status IN ('open', 'finalized')),
    UNIQUE (story_id, round_number)
);

CREATE INDEX estimation_rounds_project_story ON estimation_rounds (project_id, story_id);

CREATE TABLE estimation_round_participants (
    round_id TEXT NOT NULL REFERENCES estimation_rounds(id),
    member_id TEXT NOT NULL REFERENCES members(id),
    PRIMARY KEY (round_id, member_id)
);

CREATE TABLE estimation_votes (
    round_id TEXT NOT NULL,
    member_id TEXT NOT NULL,
    story_points INTEGER NOT NULL CHECK (story_points >= 0),
    PRIMARY KEY (round_id, member_id),
    FOREIGN KEY (round_id, member_id)
        REFERENCES estimation_round_participants(round_id, member_id)
);
