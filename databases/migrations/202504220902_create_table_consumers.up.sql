CREATE TABLE consumers (
  id VARCHAR(255) PRIMARY KEY NOT NULL,
  NIK VARCHAR(20) NOT NULL,
  full_name VARCHAR(255) NOT NULL,
  legal_name VARCHAR(255) NOT NULL,
  place_of_birth VARCHAR(255) NOT NULL,
  date_of_birth TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  salary DECIMAL(12,2) UNSIGNED NOT NULL,
  ktp_img_url VARCHAR(255) NOT NULL,
  selfie_img_url VARCHAR(255) NOT NULL,

  status_id VARCHAR(255) DEFAULT "1",
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  created_by VARCHAR(255) NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_by VARCHAR(255) NULL,
  INDEX index_NIK (NIK),
  INDEX index_full_name (full_name),
  INDEX index_legal_name (legal_name),
  INDEX index_place_of_birth (place_of_birth),
  INDEX index_date_of_birth (date_of_birth)
);