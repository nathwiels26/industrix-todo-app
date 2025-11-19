-- Drop todos table and indexes
DROP INDEX IF EXISTS idx_todos_title_search;
DROP INDEX IF EXISTS idx_todos_created_at;
DROP INDEX IF EXISTS idx_todos_category_id;
DROP INDEX IF EXISTS idx_todos_priority;
DROP INDEX IF EXISTS idx_todos_completed;
DROP TABLE IF EXISTS todos;
