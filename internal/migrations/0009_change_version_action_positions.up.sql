alter table version_actions drop column if exists "start";
alter table version_actions drop column if exists "dest";

ALTER TABLE version_actions ADD COLUMN "pos" INTEGER NOT NULL;

