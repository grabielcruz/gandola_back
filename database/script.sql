\c luis
DROP DATABASE gandola_soft;

CREATE DATABASE gandola_soft;
\c gandola_soft;

CREATE TYPE transaction_type AS ENUM ('output', 'input');

CREATE TABLE transactions (
  id SERIAL PRIMARY KEY,
  type transaction_type NOT NULL,
  amount DECIMAL CHECK (amount > 0) NOT NULL,
  description TEXT,
  executed TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  prev_balance INT
);

CREATE TABLE balances (
  id SERIAL PRIMARY KEY,
  balance DECIMAL CHECK (balance > 0) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  prev_transaction INT
);

ALTER TABLE transactions ADD CONSTRAINT fk_prev_balance FOREIGN KEY (prev_balance) REFERENCES balances(id);
ALTER TABLE balances ADD CONSTRAINT fk_prev_transaction FOREIGN KEY (prev_transaction) REFERENCES transactions(id);

-- INSERT INTO transactions(type, amount, description) 
--   VALUES ('input', '55.23', 'grasa');

-- INSERT INTO transactions(type, amount, description) 
--   VALUES ('output', '23.00', 'bornes');