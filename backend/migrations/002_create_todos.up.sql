-- Create todos table
CREATE TABLE IF NOT EXISTS todos (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE,
    priority VARCHAR(20) DEFAULT 'medium',
    due_date TIMESTAMP WITH TIME ZONE,
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for common queries
CREATE INDEX idx_todos_completed ON todos(completed);
CREATE INDEX idx_todos_priority ON todos(priority);
CREATE INDEX idx_todos_category_id ON todos(category_id);
CREATE INDEX idx_todos_created_at ON todos(created_at DESC);
CREATE INDEX idx_todos_title_search ON todos USING gin(to_tsvector('english', title));

-- Insert sample todos
INSERT INTO todos (title, description, priority, category_id, completed) VALUES
    ('Complete coding challenge', 'Build a full-stack todo application for Industrix', 'high', 1, false),
    ('Review React documentation', 'Study React hooks and context API', 'medium', 1, false),
    ('Buy groceries', 'Milk, eggs, bread, vegetables', 'low', 3, false),
    ('Morning exercise', 'Go for a 30-minute run', 'medium', 4, true),
    ('Call mom', 'Weekly catch-up call', 'high', 2, false);
