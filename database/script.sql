\c luis
DROP DATABASE gandola_soft;

CREATE DATABASE gandola_soft;
\c gandola_soft;

CREATE TYPE transaction_type AS ENUM ('output', 'input');
CREATE TYPE currency_type AS ENUM('USD', 'VES');
CREATE TYPE urgency_type AS ENUM('low', 'medium', 'high', 'critical');
CREATE TYPE actor_type AS ENUM('personnel', 'third', 'mine', 'contractee', 'driver');
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
INSERT INTO actors (type, name, national_id, address, notes) VALUES ('contractee', 'Compañía cero', 'no id', 'no address', 'no notes');
INSERT INTO actors (type, name, national_id, address, notes) VALUES ('driver', 'Conductor cero', 'no id', 'no address', 'no notes');

CREATE TABLE bills (
  id SERIAL PRIMARY KEY,
  code TEXT NOT NULL,
  url TEXT NOT NULL,
  date DATE DEFAULT CURRENT_DATE,
  company INT REFERENCES actors(id) ON DELETE RESTRICT NOT NULL,
  charged BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO bills (code, url, company) VALUES (1, 'url', 2);

CREATE TABLE trucks (
  id SERIAL PRIMARY KEY,
  name TEXT UNIQUE NOT NULL,
  data TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO trucks (name, data) VALUES ('primer camion', 'bla bla \n bla bla bla');

CREATE TABLE truck_docs (
  id SERIAL PRIMARY KEY,
  truck INT REFERENCES trucks(id) ON DELETE RESTRICT NOT NULL,
  url TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE trips (
  id SERIAL PRIMARY KEY,
  date DATE DEFAULT CURRENT_DATE,
  origin INT REFERENCES actors(id) ON DELETE RESTRICT NOT NULL,
  destination INT REFERENCES actors(id) ON DELETE RESTRICT NOT NULL,
  cargo TEXT NOT NULL,
  amount INT NOT NULL,
  unit TEXT NOT NULL,
  driver INT REFERENCES actors(id) ON DELETE RESTRICT NOT NULL,
  truck INT REFERENCES trucks(id) ON DELETE RESTRICT NOT NULL,
  bill INT REFERENCES bills(id) ON DELETE RESTRICT,
  voucher_url TEXT,
  complete BOOLEAN NOT NULL DEFAULT FALSE,
  notes TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO trips (origin, destination, cargo, amount, unit, driver, truck, voucher_url, notes) VALUES (1, 1, 'piedra', 25, 'metros', 3, 1, 'no_image', 'notes');

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