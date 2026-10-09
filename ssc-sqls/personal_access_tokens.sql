-- Personal access tokens for Super Supply Chain.
-- Store only the SHA-256 hex digest in token_hash. Never insert a plaintext token.
-- The API process also AutoMigrates this table on startup. Apply this script when
-- the database user cannot create tables.

USE super_supply_chain;

CREATE TABLE IF NOT EXISTS personal_access_tokens (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  `name` VARCHAR(100) NOT NULL,
  token_hash CHAR(64) NOT NULL,
  `prefix` VARCHAR(24) NOT NULL,
  scopes VARCHAR(500) NOT NULL DEFAULT '',
  created_at DATETIME(3) NULL,
  last_used_at DATETIME(3) NULL,
  expires_at DATETIME(3) NULL,
  revoked_at DATETIME(3) NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_personal_access_tokens_token_hash (token_hash),
  KEY idx_personal_access_tokens_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
