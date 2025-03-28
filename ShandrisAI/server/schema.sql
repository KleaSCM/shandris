-- Create the database (Note: This must be run separately in PostgreSQL)
-- CREATE DATABASE shandris_ai;

-- Enable the uuid-ossp extension if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create system_memory table
CREATE TABLE IF NOT EXISTS system_memory (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create long_term_memory table
CREATE TABLE IF NOT EXISTS long_term_memory (
    session_id TEXT NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (session_id, key)
);

-- Create persona_memory table
CREATE TABLE IF NOT EXISTS persona_memory (
    session_id TEXT PRIMARY KEY,
    traits JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create chat_history table
CREATE TABLE IF NOT EXISTS chat_history (
    id SERIAL PRIMARY KEY,
    session_id TEXT NOT NULL,
    user_message TEXT NOT NULL,
    ai_response TEXT NOT NULL,
    topic TEXT NOT NULL DEFAULT 'uncategorized',
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create session_context table
CREATE TABLE IF NOT EXISTS session_context (
    session_id TEXT PRIMARY KEY,
    current_topic TEXT NOT NULL DEFAULT 'uncategorized',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create persona_profiles table
CREATE TABLE IF NOT EXISTS persona_profiles (
    session_id TEXT PRIMARY KEY,
    profile_data JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add updated_at columns if they don't exist
DO $$ 
BEGIN 
    BEGIN
        ALTER TABLE system_memory ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
    EXCEPTION WHEN OTHERS THEN
        NULL;
    END;

    BEGIN
        ALTER TABLE long_term_memory ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
    EXCEPTION WHEN OTHERS THEN
        NULL;
    END;

    BEGIN
        ALTER TABLE persona_memory ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
    EXCEPTION WHEN OTHERS THEN
        NULL;
    END;

    BEGIN
        ALTER TABLE session_context ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
    EXCEPTION WHEN OTHERS THEN
        NULL;
    END;

    BEGIN
        ALTER TABLE persona_profiles ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
    EXCEPTION WHEN OTHERS THEN
        NULL;
    END;
END $$;

-- Add created_at columns if they don't exist
DO $$ 
BEGIN 
    BEGIN
        ALTER TABLE system_memory ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
    EXCEPTION WHEN OTHERS THEN
        NULL;
    END;

    BEGIN
        ALTER TABLE long_term_memory ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
    EXCEPTION WHEN OTHERS THEN
        NULL;
    END;

    BEGIN
        ALTER TABLE persona_memory ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
    EXCEPTION WHEN OTHERS THEN
        NULL;
    END;

    BEGIN
        ALTER TABLE session_context ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
    EXCEPTION WHEN OTHERS THEN
        NULL;
    END;

    BEGIN
        ALTER TABLE persona_profiles ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
    EXCEPTION WHEN OTHERS THEN
        NULL;
    END;
END $$;

-- Create update trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Drop existing triggers if they exist
DROP TRIGGER IF EXISTS update_system_memory_updated_at ON system_memory;
DROP TRIGGER IF EXISTS update_long_term_memory_updated_at ON long_term_memory;
DROP TRIGGER IF EXISTS update_persona_memory_updated_at ON persona_memory;
DROP TRIGGER IF EXISTS update_session_context_updated_at ON session_context;
DROP TRIGGER IF EXISTS update_persona_profiles_updated_at ON persona_profiles;

-- Create update triggers for all tables
CREATE TRIGGER update_system_memory_updated_at
    BEFORE UPDATE ON system_memory
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_long_term_memory_updated_at
    BEFORE UPDATE ON long_term_memory
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_persona_memory_updated_at
    BEFORE UPDATE ON persona_memory
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_session_context_updated_at
    BEFORE UPDATE ON session_context
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_persona_profiles_updated_at
    BEFORE UPDATE ON persona_profiles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert default AI name
INSERT INTO system_memory (key, value)
VALUES ('ai_name', 'Shandris')
ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value;

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_chat_history_session_topic ON chat_history(session_id, topic);
CREATE INDEX IF NOT EXISTS idx_long_term_memory_session ON long_term_memory(session_id);
CREATE INDEX IF NOT EXISTS idx_persona_profiles_name ON persona_profiles((profile_data->>'name')); 