DROP SCHEMA public CASCADE;
CREATE SCHEMA public;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE chat_type AS ENUM ('dialog', 'group', 'channel');
CREATE TYPE user_type AS ENUM ('owner', 'member');
CREATE TYPE reaction_type AS ENUM ('like', 'dislike');

CREATE TABLE IF NOT EXISTS public.user (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    avatar_path TEXT CHECK (avatar_path IS NULL OR (LENGTH(avatar_path) > 0 AND LENGTH(avatar_path) <= 255)),
    first_name TEXT CHECK (LENGTH(first_name) > 0 AND LENGTH(first_name) <= 50),
    last_name TEXT CHECK (LENGTH(last_name) > 0 AND LENGTH(last_name) <= 50),
    username TEXT UNIQUE NOT NULL CHECK (LENGTH(username) >= 3 AND LENGTH(username) <= 30 AND username ~ '^[a-zA-Z0-9_]+$'),
    phone TEXT UNIQUE NOT NULL CHECK (phone ~ '^[0-9]{10,15}$'),
    email TEXT UNIQUE CHECK (email IS NULL OR (LENGTH(email) <= 255 AND email ~ '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$')),
    password TEXT NOT NULL CHECK (LENGTH(password) >= 8 AND LENGTH(password) <= 72),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP CHECK (created_at <= CURRENT_TIMESTAMP),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP CHECK (updated_at >= created_at AND updated_at <= CURRENT_TIMESTAMP)
);

CREATE TABLE IF NOT EXISTS public.chat (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    avatar_path TEXT CHECK (avatar_path IS NULL OR (LENGTH(avatar_path) > 0 AND LENGTH(avatar_path) <= 255)),
    type chat_type NOT NULL,
    title TEXT NOT NULL CHECK (LENGTH(title) > 0 AND LENGTH(title) <= 100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP CHECK (created_at <= CURRENT_TIMESTAMP),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP CHECK (updated_at >= created_at AND updated_at <= CURRENT_TIMESTAMP)
);

CREATE TABLE IF NOT EXISTS public.user_chat (
    user_id UUID NOT NULL,
    chat_id UUID NOT NULL,
    user_role user_type NOT NULL,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP CHECK (joined_at <= CURRENT_TIMESTAMP),
    PRIMARY KEY (user_id, chat_id),
    FOREIGN KEY (user_id) REFERENCES public.user(id) ON DELETE CASCADE,
    FOREIGN KEY (chat_id) REFERENCES public.chat(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS public.contact (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    contact_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP CHECK (created_at <= CURRENT_TIMESTAMP),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP CHECK (updated_at >= created_at AND updated_at <= CURRENT_TIMESTAMP),
    CHECK (user_id <> contact_id),
    UNIQUE (user_id, contact_id),
    FOREIGN KEY (user_id) REFERENCES public.user(id) ON DELETE CASCADE,
    FOREIGN KEY (contact_id) REFERENCES public.user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS public.message (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_message_id UUID,
    chat_id UUID NOT NULL,
    user_id UUID NOT NULL,
    body TEXT NOT NULL CHECK (LENGTH(body) > 0 AND LENGTH(body) <= 2000),
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP CHECK (sent_at <= CURRENT_TIMESTAMP),
    is_redacted BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (parent_message_id) REFERENCES public.message(id) ON DELETE SET NULL,
    FOREIGN KEY (chat_id) REFERENCES public.chat(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES public.user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS public.message_reaction (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id UUID NOT NULL,
    user_id UUID NOT NULL,
    reaction reaction_type NOT NULL,
    reacted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP CHECK (reacted_at <= CURRENT_TIMESTAMP),
    UNIQUE (message_id, user_id),
    FOREIGN KEY (message_id) REFERENCES public.message(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES public.user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS public.message_view (
    message_id UUID NOT NULL,
    user_id UUID NOT NULL,
    viewed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP CHECK (viewed_at <= CURRENT_TIMESTAMP),
    PRIMARY KEY (message_id, user_id),
    FOREIGN KEY (message_id) REFERENCES public.message(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES public.user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS public.message_payload (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id UUID NOT NULL,
    file_path TEXT NOT NULL CHECK (LENGTH(file_path) > 0 AND LENGTH(file_path) <= 255),
    FOREIGN KEY (message_id) REFERENCES public.message(id) ON DELETE CASCADE
);