alter table repositories drop column if exists filename;

alter table repositories drop constraint oma_repo_id_unique_for_each_filename

