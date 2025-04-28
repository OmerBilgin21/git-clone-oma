alter table version_actions drop column if exists start_x;
alter table version_actions drop column if exists start_y;
alter table version_actions drop column if exists end_x;
alter table version_actions drop column if exists end_y;

alter table version_actions add column "start" INTEGER default null;
alter table version_actions add column "dest" INTEGER not null;

