\c luis
DROP DATABASE gandola_soft;

CREATE DATABASE gandola_soft;
\c gandola_soft;

CREATE TYPE transaction_type AS ENUM ('output', 'input');
CREATE TYPE currency_type AS ENUM('USD', 'VES');
CREATE TYPE urgency_type AS ENUM('low', 'medium', 'high', 'critical');
CREATE TYPE actor_type AS ENUM('personnel', 'third', 'mine', 'contractee');
CREATE EXTENSION CITEXT;
-- tipos de actores:
--   - El empleado: Luis D, papa, yo, Niliberto
--   - El tercero: Mr frenos, toro mocho, ochoa, simpson
--   - El saque: San Remo, Farias
--   - El contratante: Cayucos, Nivar, Super S, Proporca, Bicolor

CREATE TABLE actors (
  id SERIAL PRIMARY KEY,
  type actor_type NOT NULL,
  name CITEXT NOT NULL UNIQUE,
  national_id TEXT,
  address TEXT,
  notes TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO actors (type, name, national_id, address, notes) VALUES ('third', 'Externo', 'no id', 'no address', 'no notes');

CREATE TABLE bills (
  id SERIAL PRIMARY KEY,
  url TEXT NOT NULL,
  date TIME WITH TIME ZONE DEFAULT CURRENT_TIME,
  company INT REFERENCES actors(id) ON DELETE RESTRICT NOT NULL,
  charged BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO bills (url, company) VALUES ('url', 1);

CREATE TABLE trips (
  id SERIAL PRIMARY KEY,
  date TIME WITH TIME ZONE DEFAULT CURRENT_TIME,
  origin INT REFERENCES actors(id) ON DELETE RESTRICT NOT NULL,
  destination INT REFERENCES actors(id) ON DELETE RESTRICT NOT NULL,
  cargo TEXT NOT NULL,
  driver TEXT NOT NULL,
  truck TEXT NOT NULL,
  bill INT REFERENCES bills(id) ON DELETE RESTRICT,
  voucher TEXT,
  complete BOOLEAN NOT NULL DEFAULT FALSE,
  notes TEXT 
);

INSERT INTO trips (origin, destination, cargo, driver, truck, voucher, notes) VALUES (1, 1, '25 metros piedra', 'chofer', 'camion', 'vauche', 'notes' );

CREATE TABLE transactions_with_balances (
  id SERIAL PRIMARY KEY,
  type transaction_type NOT NULL,
  currency currency_type NOT NULL,
  amount DECIMAL(17,2) CHECK (amount >= 0) NOT NULL,
  description TEXT NOT NULL,
  USD_balance DECIMAL(22,2) CHECK (USD_balance >= 0) NOT NULL,
  VES_balance DECIMAL(22,2) CHECK (VES_balance >= 0) NOT NULL,
  actor INT REFERENCES actors(id) ON DELETE RESTRICT NOT NULL,
  executed TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO transactions_with_balances (type, currency, amount, description, USD_balance, VES_balance, actor)
  VALUES ('input', 'USD', '0', 'transaction zero', '0', '0', '1');

CREATE TABLE pending_transactions (
  id SERIAL PRIMARY KEY,
  type transaction_type NOT NULL,
  currency currency_type NOT NULL,
  amount DECIMAL(17,2) CHECK (amount >= 0) NOT NULL,
  description TEXT NOT NULL,
  actor INT REFERENCES actors(id) ON DELETE RESTRICT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO pending_transactions (type, currency, amount, description, actor) 
  VALUES ('input', 'USD', '0', 'pending transaction zero', '1');

CREATE TABLE notes (
  id SERIAL PRIMARY KEY,
  description TEXT NOT NULL,
  urgency urgency_type NOT NULL,
  attended BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  attended_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO notes (description, urgency) VALUES ('first note', 'low');

-- CREATE TABLE docs (
--   id SERIAL PRIMARY KEY,
--   description TEXT NOT NULL,
--   url TEXT NOT NULL,
--   created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
-- );