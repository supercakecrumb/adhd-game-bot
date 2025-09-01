-- Migration 005: Add timezone support to tasks
BEGIN;

ALTER TABLE tasks
ADD COLUMN time_zone VARCHAR(50) NOT NULL DEFAULT 'UTC';

COMMIT;