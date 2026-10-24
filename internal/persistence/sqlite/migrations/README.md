# SQLite migrations

Add versioned SQL files here using the format `NNNNNN_description.sql`.
The migration runner embeds this directory and applies files in lexical order.
The baseline migration is intentionally empty. `000002_projects.sql` creates the
shared project table used by US-01 and dependent stories. `000003_backlog_items.sql`
creates Product Backlog items. `000004_sprints.sql` creates Sprint records for
US-07, linked to existing projects by their UUID. `000005_sprint_backlog_items.sql`
stores immutable planning-time snapshots for US-08.

`000006_members.sql` and `000007_members_name_storage.sql` add project members
and preserve exact accepted names. `000008_projects_dates.sql` adds nullable
calendar-day start/end columns for US-02 without changing existing rows or
relations. Dates are not inferred or filled with the migration date.

`000009_sprint_completion.sql` records verified completion separately from the
immutable US-08 planning status.

`000010_sprint_closure.sql` persists whether a Sprint is closed for US-10.

`000011_planning_poker.sql` adds the initial estimation round, its participant
snapshot, and one replaceable vote per participant for US-12.
