-- US-02: preserve existing values and relations, without inventing past dates.
ALTER TABLE projects ADD COLUMN start_date TEXT NULL;
ALTER TABLE projects ADD COLUMN end_date TEXT NULL;
