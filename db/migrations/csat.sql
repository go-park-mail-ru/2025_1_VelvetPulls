-- Удаляем схему и ее содержимое, если она существует
DROP SCHEMA IF EXISTS csat CASCADE;
CREATE SCHEMA csat;

-- Устанавливаем расширение pgcrypto
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Создаем таблицу question с вопросами
CREATE TABLE IF NOT EXISTS csat.question (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL CONSTRAINT title_not_empty CHECK (length(title) > 0),
    question_text TEXT NOT NULL CONSTRAINT question_text_not_empty CHECK (length(question_text) > 0),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(title)
);

-- Создаем таблицу answer с рейтингами
CREATE TABLE IF NOT EXISTS csat.answer (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_id UUID NOT NULL,
    username TEXT NOT NULL,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),  -- Рейтинг от 1 до 5
    feedback TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP CHECK (created_at <= CURRENT_TIMESTAMP),
    FOREIGN KEY (question_id) REFERENCES csat.question(id) ON DELETE CASCADE
);

-- Создаем таблицу user_activity для отслеживания активности пользователей
CREATE TABLE IF NOT EXISTS csat.user_activity (
    username TEXT UNIQUE NOT NULL,
    last_response_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    responses_count INTEGER DEFAULT 0,
    PRIMARY KEY (username)
);

-- Добавляем вопрос в таблицу question
INSERT INTO csat.question (title, question_text)
VALUES
('Создание диалога', 'Удобно ли начать диалог?'),
('Создание чата', 'Удобно ли создать чат?'),
('Про диалог', 'Нравится ли общение в диалогах?'),
('Про группы', 'Удобны ли групповые чаты?'),
('Про профиль', 'Удобно ли настроить профиль?'),
('Создание контактов', 'Удобно ли добавлять контакты?'),
('Просмотр контактов', 'Удобно ли просматривать контакты?'),
('Регистрация', 'Удобно ли регистрироваться?');
