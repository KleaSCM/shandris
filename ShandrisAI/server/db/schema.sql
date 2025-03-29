-- Core system tables
CREATE TABLE IF NOT EXISTS moods (
    id UUID PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    current_value FLOAT NOT NULL,
    base_value FLOAT NOT NULL,
    last_updated TIMESTAMP NOT NULL,
    decay_rate FLOAT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS mood_patterns (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    keywords TEXT[] NOT NULL,
    sentiment FLOAT NOT NULL,
    mood_shift VARCHAR(50) NOT NULL,
    intensity FLOAT NOT NULL,
    decay_rate FLOAT NOT NULL,
    requirements TEXT[],
    exclusions TEXT[],
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS traits (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    trait_name VARCHAR(100) NOT NULL,
    value FLOAT NOT NULL,
    confidence FLOAT NOT NULL,
    last_updated TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS topics (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    keywords TEXT[] NOT NULL,
    last_discussed TIMESTAMP,
    frequency INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Timeline and Memory tables
CREATE TABLE IF NOT EXISTS memory_events (
    id UUID PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    importance FLOAT NOT NULL,
    context JSONB NOT NULL,
    relations TEXT[] NOT NULL,
    tags TEXT[] NOT NULL,
    emotions JSONB NOT NULL,
    last_recall TIMESTAMP,
    recall_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS timeline_markers (
    id UUID PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    recurrence_pattern JSONB,
    importance FLOAT NOT NULL,
    last_trigger TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Persona System tables
CREATE TABLE IF NOT EXISTS personas (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    traits JSONB NOT NULL,
    mood_bias JSONB NOT NULL,
    style_rules JSONB NOT NULL,
    preferences JSONB NOT NULL,
    constraints JSONB NOT NULL,
    active BOOLEAN NOT NULL DEFAULT false,
    last_used TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS persona_transitions (
    id UUID PRIMARY KEY,
    from_persona UUID REFERENCES personas(id),
    to_persona UUID REFERENCES personas(id),
    timestamp TIMESTAMP NOT NULL,
    reason TEXT NOT NULL,
    context JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Topic Memory and Relationships tables
CREATE TABLE IF NOT EXISTS topic_relationships (
    id UUID PRIMARY KEY,
    from_topic UUID REFERENCES topics(id),
    to_topic UUID REFERENCES topics(id),
    relationship_type VARCHAR(50) NOT NULL,
    strength FLOAT NOT NULL,
    bidirectional BOOLEAN NOT NULL,
    context JSONB NOT NULL,
    last_active TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS facts (
    id UUID PRIMARY KEY,
    content TEXT NOT NULL,
    source VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    confidence FLOAT NOT NULL,
    relations TEXT[] NOT NULL,
    context JSONB NOT NULL,
    verification JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Session Management tables
CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    start_time TIMESTAMP NOT NULL,
    last_active TIMESTAMP NOT NULL,
    current_state JSONB NOT NULL,
    context JSONB NOT NULL,
    active_systems JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS session_checkpoints (
    id UUID PRIMARY KEY,
    session_id UUID REFERENCES sessions(id),
    timestamp TIMESTAMP NOT NULL,
    state JSONB NOT NULL,
    context JSONB NOT NULL,
    metadata JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_moods_name ON moods(name);
CREATE INDEX IF NOT EXISTS idx_traits_user_id ON traits(user_id);
CREATE INDEX IF NOT EXISTS idx_topics_category ON topics(category);
CREATE INDEX IF NOT EXISTS idx_memory_events_type ON memory_events(type);
CREATE INDEX IF NOT EXISTS idx_memory_events_timestamp ON memory_events(timestamp);
CREATE INDEX IF NOT EXISTS idx_timeline_markers_type ON timeline_markers(type);
CREATE INDEX IF NOT EXISTS idx_timeline_markers_timestamp ON timeline_markers(timestamp);
CREATE INDEX IF NOT EXISTS idx_personas_type ON personas(type);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_last_active ON sessions(last_active);

-- Add GiST index for text search on topics
CREATE INDEX IF NOT EXISTS idx_topics_keywords ON topics USING GIN (keywords);

-- Add JSON indexes
CREATE INDEX IF NOT EXISTS idx_memory_events_context ON memory_events USING GIN (context);
CREATE INDEX IF NOT EXISTS idx_memory_events_emotions ON memory_events USING GIN (emotions);
CREATE INDEX IF NOT EXISTS idx_personas_traits ON personas USING GIN (traits);
CREATE INDEX IF NOT EXISTS idx_sessions_state ON sessions USING GIN (current_state); 