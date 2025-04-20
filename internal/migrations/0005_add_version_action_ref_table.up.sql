CREATE TABLE version_actions (
  id SERIAL PRIMARY KEY,
  start_x INTEGER NOT NULL,
  start_y INTEGER NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP WITHOUT TIME ZONE,
  end_x INTEGER NOT NULL,
  end_y INTEGER NOT NULL,
  action_key TEXT NOT NULL CHECK (action_key IN ('addition', 'deletion')),
  version_id INTEGER NOT NULL
);

ALTER TABLE version_actions
ADD CONSTRAINT fk_version
FOREIGN KEY (version_id)
REFERENCES versions(id)
ON DELETE CASCADE;

