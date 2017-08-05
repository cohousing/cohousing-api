-- +migrate Up
CREATE TABLE users
(
  id           INT(10) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  created_at   TIMESTAMP   NOT NULL,
  updated_at   TIMESTAMP   NOT NULL,
  deleted_at   TIMESTAMP   NULL,
  username     VARCHAR(64) NOT NULL UNIQUE,
  password     VARCHAR(60) NOT NULL,
  global_admin TINYINT(1)  NOT NULL,
  residents    CHAR(4)     NOT NULL,
  apartments   CHAR(4)     NOT NULL
)
  ENGINE = InnoDB;

CREATE TABLE groups
(
  id           INT(10) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  created_at   TIMESTAMP    NOT NULL,
  updated_at   TIMESTAMP    NOT NULL,
  deleted_at   TIMESTAMP    NULL,
  name         VARCHAR(100) NOT NULL UNIQUE,
  global_admin TINYINT(1)   NOT NULL,
  residents    CHAR(4)      NOT NULL,
  apartments   CHAR(4)      NOT NULL
)
  ENGINE = InnoDB;

CREATE TABLE users_groups
(
  user_id  INT(10) UNSIGNED NOT NULL,
  group_id INT(10) UNSIGNED NOT NULL,
  PRIMARY KEY (user_id, group_id),
  CONSTRAINT fk_user_groups_user FOREIGN KEY (user_id) REFERENCES users (id),
  CONSTRAINT fk_user_groups_group FOREIGN KEY (group_id) REFERENCES groups (id)
)
  ENGINE = InnoDB;

-- +migrate Down
DROP TABLE users_groups;
DROP TABLE users;
DROP TABLE groups;