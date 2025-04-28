alter table version_actions add column start_x INTEGER NOT NULL;
alter table version_actions add column start_y INTEGER NOT NULL;
alter table version_actions add column end_x INTEGER NOT NULL;
alter table version_actions add column end_y INTEGER NOT NULL;

alter table version_actions drop column if exists start;
alter table version_actions drop column if exists dest;

