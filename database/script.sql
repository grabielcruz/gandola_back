\c luis
DROP DATABASE gandola_soft;

CREATE DATABASE gandola_soft;
\c gandola_soft;

CREATE TYPE transaction_type AS ENUM ('output', 'input');

CREATE TABLE actors (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  description TEXT NOT NULL,
  is_company BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO actors (name, description) VALUES ('Externo', 'renglón para actor no registrado');

CREATE TABLE transactions_with_balances (
  id SERIAL PRIMARY KEY,
  type transaction_type NOT NULL,
  amount DECIMAL(9, 2) CHECK (amount >= 0) NOT NULL,
  description TEXT NOT NULL,
  balance DECIMAL(9, 2) CHECK (balance >= 0) NOT NULL,
  actor INT REFERENCES actors(id) ON DELETE RESTRICT NOT NULL,
  executed TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO transactions_with_balances (type, amount, description, balance, actor)
  VALUES ('input', '0', 'transaction zero', '0', '1');

CREATE TABLE pending_transactions (
  id SERIAL PRIMARY KEY,
  type transaction_type NOT NULL,
  amount DECIMAL(9, 2) CHECK (amount >= 0) NOT NULL,
  description TEXT NOT NULL,
  actor INT REFERENCES actors(id) ON DELETE RESTRICT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO pending_transactions (type, amount, description, actor) 
  VALUES ('input', '0', 'pending transaction zero', '1');

CREATE TABLE trip_bills (
  id TEXT PRIMARY KEY,
  url TEXT NOT NULL,
  date TIME WITH TIME ZONE DEFAULT CURRENT_TIME,
  company INT REFERENCES actors(id) ON DELETE RESTRICT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE trips (
  id SERIAL PRIMARY KEY,
  date TIME WITH TIME ZONE DEFAULT CURRENT_TIME,
  origin TEXT NOT NULL,
  destination TEXT NOT NULL,
  cargo TEXT NOT NULL,
  driver TEXT NOT NULL,
  truck TEXT NOT NULL,
  bill TEXT REFERENCES trip_bills(id) ON DELETE RESTRICT,
  support TEXT,
  notes TEXT 
);

CREATE TABLE docs (
  id SERIAL PRIMARY KEY,
  description TEXT NOT NULL,
  url TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);