\c luis
DROP DATABASE gandola_soft;

CREATE DATABASE gandola_soft;
\c gandola_soft;

-- CREATE TYPE transaction_type AS ENUM ('output', 'input');

CREATE TABLE transactions (
  id SERIAL PRIMARY KEY,
  -- type transaction_type,
  amount DECIMAL,
  executed TIMESTAMP
);

INSERT INTO transactions(amount, executed) 
  VALUES ('55.23', '2021-06-17 04:05:06 -4:00');

INSERT INTO transactions(amount, executed) 
  VALUES ('-23.00', '2021-05-09 08:09:12 -4:00');