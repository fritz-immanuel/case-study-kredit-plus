CREATE TABLE consumer_transactions (
  id VARCHAR(255) PRIMARY KEY NOT NULL,
  consumer_id VARCHAR(255) NOT NULL,
  contract_number VARCHAR(255) NOT NULL,
  OTR DECIMAL(12,2) UNSIGNED NOT NULL,
  admin_fee DECIMAL(12,2) UNSIGNED NOT NULL,
  installment_amount DECIMAL(12,2) UNSIGNED NOT NULL,
  loan_term INT UNSIGNED NOT NULL,
  interest_amount DECIMAL(12,2) UNSIGNED NOT NULL,
  total_amount DECIMAL(12,2) UNSIGNED NOT NULL,
  asset_name VARCHAR(255) NOT NULL,

  status_id VARCHAR(255) DEFAULT "1",
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  created_by VARCHAR(255) NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_by VARCHAR(255) NULL,
  INDEX index_consumer_id (consumer_id),
  INDEX index_contract_number (contract_number)
);