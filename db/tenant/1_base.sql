-- +migrate Up
CREATE TABLE apartments
(
  id         INT(10) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  created_at TIMESTAMP    NULL,
  updated_at TIMESTAMP    NULL,
  deleted_at TIMESTAMP    NULL,
  address    VARCHAR(255) NULL
);

ALTER TABLE apartments
  ADD INDEX idx_apartments_deleted_at (deleted_at);

CREATE TABLE residents
(
  id            INT(10) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  created_at    TIMESTAMP        NULL,
  updated_at    TIMESTAMP        NULL,
  deleted_at    TIMESTAMP        NULL,
  name          VARCHAR(255)     NULL,
  phone_number  VARCHAR(255)     NULL,
  email_address VARCHAR(255)     NULL,
  apartment_id  INT(10) UNSIGNED NULL
);

ALTER TABLE residents
  ADD INDEX idx_residents_deleted_at (deleted_at),
  ADD INDEX idx_residents_apartment_id (apartment_id),
  ADD CONSTRAINT fk_residents_apartment_id FOREIGN KEY (apartment_id) REFERENCES apartments (id);

-- +migrate Down
DROP TABLE IF EXISTS residents;
DROP TABLE IF EXISTS apartments;