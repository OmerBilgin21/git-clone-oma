ALTER TABLE version_actions DROP COLUMN if exists "pos";

alter table version_actions ADD column "start" INTEGER NOT NULL;
alter table version_actions ADD column "dest" INTEGER NOT NULL;

