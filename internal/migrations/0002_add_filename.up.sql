alter table repositories add column filename varchar not null;

ALTER TABLE repositories
ADD CONSTRAINT
oma_repo_id_unique_for_each_filename
UNIQUE(oma_repo_id, filename)

