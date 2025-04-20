alter table repositories drop column if exists created_at;
alter table repositories add column created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP;

