CREATE DATABASE IF NOT EXISTS case-study-kredit-plus;

CREATE TABLE IF NOT EXISTS schema_migrations (
  version BIGINT DEFAULT 0,
  dirty VARCHAR(5) DEFAULT "",
  PRIMARY KEY (version)
);

TRUNCATE schema_migrations;
INSERT INTO schema_migrations (version, dirty) VALUES (202406161200, "0");

CREATE TABLE IF NOT EXISTS status (
  id VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  PRIMARY KEY (id)
);