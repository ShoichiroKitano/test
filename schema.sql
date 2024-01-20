CREATE DATABASE IF NOT EXISTS super_shiharai_kun;

use super_shiharai_kun;
CREATE TABLE invoices (
   id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
   corporate_id BIGINT NOT NULL,
   partner_id BIGINT NOT NULL,
   issue_date DATE NOT NULL,
   payment_amount INTEGER NOT NULL,
   fee INTEGER Not NULL,
   fee_rate decimal(2, 2) UNSIGNED Not NULL,
   sales_tax_rate_id BIGINT NOT NULL,
   amount_due INTEGER NOT NULL,
   payment_due_date DATE NOT NULL,
   status TINYINT UNSIGNED NOT NULL,
   PRIMARY KEY (id)
);

CREATE INDEX idx_corporate_id_payment_due_date ON invoices (corporate_id, payment_due_date);
