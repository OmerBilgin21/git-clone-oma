CREATE TABLE versions (
  id SERIAL PRIMARY KEY,
  start_x INTEGER NOT NULL,
  start_y INTEGER NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP WITHOUT TIME ZONE,
  end_x INTEGER NOT NULL,
  end_y INTEGER NOT NULL,
  action_key TEXT NOT NULL CHECK (action_key IN ('addition', 'deletion')),
  repository_id INTEGER NOT NULL
);

ALTER TABLE versions
ADD CONSTRAINT fk_repository
FOREIGN KEY (repository_id)
REFERENCES repositories(id)
ON DELETE CASCADE;

