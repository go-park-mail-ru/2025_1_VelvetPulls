DROP SCHEMA IF EXISTS csat CASCADE;
CREATE SCHEMA csat;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE rating_scale AS ENUM ('1', '2', '3', '4', '5');

CREATE TABLE IF NOT EXISTS csat.question (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL CONSTRAINT title_not_empty CHECK (length(title) > 0),
    question_text TEXT NOT NULL CONSTRAINT question_text_not_empty CHECK (length(question_text) > 0),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(title)
);

CREATE TABLE IF NOT EXISTS csat.answer (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_id UUID NOT NULL,
    user_id UUID,
    rating rating_scale NOT NULL,
    feedback TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP CHECK (created_at <= CURRENT_TIMESTAMP),
    FOREIGN KEY (question_id) REFERENCES csat.question(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS csat.user_activity (
    username TEXT UNIQUE NOT NULL,
    last_response_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    responses_count INTEGER DEFAULT 0,
    PRIMARY KEY (username)
);