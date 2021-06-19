\c luis
DROP DATABASE gandola_soft;

CREATE DATABASE gandola_soft;
\c gandola_soft;

CREATE TYPE transaction_type AS ENUM ('output', 'input');

CREATE TABLE transactions (
  id SERIAL PRIMARY KEY,
  type transaction_type,
  amount DECIMAL,
  description TEXT,
  executed TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO transactions(type, amount, description) 
  VALUES ('input', '55.23', 'grasa');

INSERT INTO transactions(type, amount, description) 
  VALUES ('output', '23.00', 'bornes');