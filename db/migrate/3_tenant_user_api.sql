-- +migrate Up
ALTER TABLE
users
  CHANGE COLUMN residents perm_residents CHAR(4) NOT NULL,
  CHANGE COLUMN apartments perm_apartments CHAR(4) NOT NULL,
  ADD COLUMN perm_users CHAR(4) NOT NULL DEFAULT '____';

ALTER TABLE
groups
  CHANGE COLUMN residents perm_residents CHAR(4) NOT NULL,
  CHANGE COLUMN apartments perm_apartments CHAR(4) NOT NULL,
  ADD COLUMN perm_users CHAR(4) NOT NULL DEFAULT '____';

-- +migrate Down
ALTER TABLE
users
  CHANGE COLUMN perm_residents residents CHAR(4) NOT NULL,
  CHANGE COLUMN perm_apartments apartments CHAR(4) NOT NULL,
  DROP COLUMN perm_users;

ALTER TABLE
groups
  CHANGE COLUMN perm_residents residents CHAR(4) NOT NULL,
  CHANGE COLUMN perm_apartments apartments CHAR(4) NOT NULL,
  DROP COLUMN perm_users;
