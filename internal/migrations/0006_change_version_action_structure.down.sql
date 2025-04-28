alter table add column start_x INTEGER NOT NULL;
alter table add column start_y INTEGER NOT NULL;
alter table add column end_x INTEGER NOT NULL;
alter table add column end_y INTEGER NOT NULL;

alter table drop column if exists start;
alter table drop column if exists dest;

