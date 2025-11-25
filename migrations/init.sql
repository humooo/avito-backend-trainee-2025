CREATE TABLE IF NOT EXISTS teams (
    name TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    team_name TEXT NOT NULL REFERENCES teams(name) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS pull_requests (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    author_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'OPEN',
    created_at TIMESTAMP DEFAULT NOW(),
    merged_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pr_reviewers (
    pr_id TEXT REFERENCES pull_requests(id) ON DELETE CASCADE,
    reviewer_id TEXT REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (pr_id, reviewer_id)
);
