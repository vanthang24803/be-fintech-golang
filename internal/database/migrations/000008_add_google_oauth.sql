-- Migration: Add Google ID to Users table
ALTER TABLE users ADD COLUMN google_id VARCHAR(255) UNIQUE;
