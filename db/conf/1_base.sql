-- +migrate Up
CREATE TABLE tenants
(
  id         INT(10) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  created_at TIMESTAMP    NOT NULL,
  updated_at TIMESTAMP    NULL,
  deleted_at TIMESTAMP    NULL,
  context    VARCHAR(100) NOT NULL,
  name       VARCHAR(255) NOT NULL,
  custom_url VARCHAR(255) NULL
);

ALTER TABLE tenants
  ADD INDEX idx_tenants_deleted_at (deleted_at);

-- +migrate Down
DROP TABLE tenants;