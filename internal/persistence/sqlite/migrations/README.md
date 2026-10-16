# SQLite migrations

Add versioned SQL files here using the format `NNNNNN_description.sql`.
The migration runner embeds this directory and applies files in lexical order.
The baseline migration is intentionally empty. `000002_projects.sql` creates the
shared project table used by US-01 and dependent stories. `000003_backlog_items.sql`
creates Product Backlog items. `000004_sprints.sql` creates Sprint records for
US-07, linked to existing projects by their UUID.
