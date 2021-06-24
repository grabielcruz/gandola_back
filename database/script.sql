\c luis
DROP DATABASE gandola_soft;

CREATE DATABASE gandola_soft;
\c gandola_soft;

CREATE TYPE transaction_type AS ENUM ('output', 'input', 'zero');

CREATE TABLE transactions_with_balances (
  id SERIAL PRIMARY KEY,
  type transaction_type NOT NULL,
  amount DECIMAL(9, 2) CHECK (amount >= 0) NOT NULL,
  description TEXT,
  balance DECIMAL(9, 2) CHECK (balance >= 0) NOT NULL,
  executed TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO transactions_with_balances(type, amount, description, balance)
  VALUES ('zero', '0', 'transaction zero', '0');