-- Migration: 001_add_category_id.sql
-- Adds category_id column to products table if it doesn't exist
BEGIN;
ALTER TABLE products ADD COLUMN IF NOT EXISTS category_id INTEGER REFERENCES categories(id);
COMMIT;
