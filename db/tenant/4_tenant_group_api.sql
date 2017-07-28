-- +migrate Up
ALTER TABLE
users
  ADD COLUMN perm_groups CHAR(4) NOT NULL DEFAULT '____';

ALTER TABLE
groups
  ADD COLUMN perm_groups CHAR(4) NOT NULL DEFAULT '____';

-- +migrate Down
ALTER TABLE
users
  DROP COLUMN perm_groups;

ALTER TABLE
groups
  DROP COLUMN perm_groups;
