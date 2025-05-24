DROP SCHEMA public CASCADE;
CREATE SCHEMA public;

CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS unaccent;

CREATE TYPE chat_type AS ENUM ('dialog', 'group', 'channel');
CREATE TYPE message_type AS ENUM ('default', 'with_payload', 'sticker');
CREATE TYPE user_type AS ENUM ('owner', 'member');
CREATE TYPE reaction_type AS ENUM ('like', 'dislike');

CREATE TABLE IF NOT EXISTS public.user (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    avatar_path TEXT CHECK (avatar_path IS NULL OR (LENGTH(avatar_path) > 0 AND LENGTH(avatar_path) <= 255)),
    name TEXT CHECK (LENGTH(name) > 0 AND LENGTH(name) <= 30),
    username TEXT UNIQUE NOT NULL CHECK (
        LENGTH(username) >= 3 AND 
        LENGTH(username) <= 30 AND 
        username ~ '^[a-zA-Z0-9!@#$%^&*()_+=\-]+$'
    ),
    password TEXT NOT NULL CHECK (LENGTH(password) >= 8 AND LENGTH(password) <= 72),
    birth_date DATE CHECK (birth_date <= CURRENT_DATE),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
);

CREATE TABLE IF NOT EXISTS public.contact (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    contact_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP CHECK (created_at <= CURRENT_TIMESTAMP),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP CHECK (updated_at >= created_at AND updated_at <= CURRENT_TIMESTAMP),
    CHECK (user_id <> contact_id),
    UNIQUE (user_id, contact_id),
    FOREIGN KEY (user_id) REFERENCES public.user(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (contact_id) REFERENCES public.user(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS public.chat (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    avatar_path TEXT CHECK (avatar_path IS NULL OR (LENGTH(avatar_path) > 0 AND LENGTH(avatar_path) <= 255)),
    type chat_type NOT NULL,
    title TEXT NOT NULL CHECK (LENGTH(title) > 0 AND LENGTH(title) <= 100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.user_chat (
    user_id UUID NOT NULL,
    chat_id UUID NOT NULL,
    user_role user_type NOT NULL,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP CHECK (joined_at <= CURRENT_TIMESTAMP),
    send_notifications boolean DEFAULT true NOT NULL,
    PRIMARY KEY (user_id, chat_id),
    FOREIGN KEY (user_id) REFERENCES public.user(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (chat_id) REFERENCES public.chat(id) ON DELETE CASCADE ON UPDATE CASCADE
);


CREATE TABLE IF NOT EXISTS public.message (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_message_id UUID,
    message_type message_type NOT NULL,
    chat_id UUID NOT NULL,
    user_id UUID NOT NULL,
    sticker_path text,
    body TEXT NOT NULL CHECK (LENGTH(body) <= 2000),
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP CHECK (sent_at <= CURRENT_TIMESTAMP),
    is_redacted BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (parent_message_id) REFERENCES public.message(id) ON DELETE SET NULL ON UPDATE CASCADE,
    FOREIGN KEY (chat_id) REFERENCES public.chat(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES public.user(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS public.message_reaction (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id UUID NOT NULL,
    user_id UUID NOT NULL,
    reaction reaction_type NOT NULL,
    reacted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP CHECK (reacted_at <= CURRENT_TIMESTAMP),
    UNIQUE (message_id, user_id),
    FOREIGN KEY (message_id) REFERENCES public.message(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES public.user(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS public.message_view (
    message_id UUID NOT NULL,
    user_id UUID NOT NULL,
    viewed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP CHECK (viewed_at <= CURRENT_TIMESTAMP),
    PRIMARY KEY (message_id, user_id),
    FOREIGN KEY (message_id) REFERENCES public.message(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES public.user(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS public.message_payload (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id UUID NOT NULL,
    file_path TEXT NOT NULL CHECK (LENGTH(file_path) > 0 AND LENGTH(file_path) <= 255),
    file_name TEXT NOT NULL CHECK (LENGTH(file_name) > 0 AND LENGTH(file_name) <= 255),
    content_type TEXT NOT NULL,
    file_size INT CHECK (file_size >= 0),
    FOREIGN KEY (message_id) REFERENCES public.message(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS public.sticker (
	id uuid NOT NULL,
	sticker_path text NOT NULL,
	CONSTRAINT sticker_pk PRIMARY KEY (id),
	CONSTRAINT sticker_unique UNIQUE (sticker_path)
);

CREATE TABLE IF NOT EXISTS public.sticker_pack (
	id uuid NOT NULL,
	"name" text NULL,
	photo_id text NOT NULL,
	CONSTRAINT sticker_pack_pk PRIMARY KEY (id)
);

CREATE TABLE public.sticker_sticker_pack (
	id uuid NOT NULL,
	sticker uuid NOT NULL,
	pack uuid NOT NULL,
	CONSTRAINT sticker_sticker_pack_pk PRIMARY KEY (id),
	CONSTRAINT sticker_sticker_pack_sticker_fk FOREIGN KEY (sticker) REFERENCES public.sticker(id),
	CONSTRAINT sticker_sticker_pack_sticker_pack_fk FOREIGN KEY (pack) REFERENCES public.sticker_pack(id)
);

DROP INDEX IF EXISTS idx_user_search_gin;
DROP INDEX IF EXISTS idx_chat_title_gin;
CREATE INDEX idx_user_username_trgm ON public.user USING gin (username gin_trgm_ops);
CREATE INDEX idx_user_name_trgm ON public.user USING gin (name gin_trgm_ops);

CREATE INDEX idx_chat_title_trgm ON chat USING gin (title gin_trgm_ops);
CREATE INDEX idx_message_body_trgm ON message USING gin (body gin_trgm_ops);

CREATE INDEX idx_message_chat_sent_at ON message(chat_id, sent_at DESC);
CREATE INDEX idx_message_user_id ON message(user_id);
