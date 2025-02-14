-- /postgres/init/01-init.sql
-- Create extensions
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Create databases
CREATE DATABASE users;
CREATE DATABASE addresses;
