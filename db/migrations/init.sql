DROP SCHEMA public CASCADE;
CREATE SCHEMA public;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE chat_type AS ENUM ('dialog', 'group', 'channel');
CREATE TYPE user_type AS ENUM ('owner', 'member');
CREATE TYPE reaction_type AS ENUM ('like', 'dislike');

CREATE TABLE IF NOT EXISTS public.user (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    avatar_path TEXT,
    first_name TEXT,
    last_name TEXT,
    username TEXT UNIQUE NOT NULL,
    phone TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.chat (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    avatar_path TEXT,
    type chat_type NOT NULL,
    title TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.user_chat (
  user_id UUID NOT NULL,
  chat_id UUID NOT NULL,
  user_role user_type NOT NULL,
  joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (user_id, chat_id),
  FOREIGN KEY (user_id) REFERENCES public.user(id) ON DELETE CASCADE,
  FOREIGN KEY (chat_id) REFERENCES public.chat(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS public.contact (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    contact_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
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
    body TEXT NOT NULL,
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_redacted BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (parent_message_id) REFERENCES public.message(id) ON DELETE CASCADE,
    FOREIGN KEY (chat_id) REFERENCES public.chat(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES public.user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS public.message_reaction (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id UUID NOT NULL,
    user_id UUID NOT NULL,
    reaction reaction_type NOT NULL,
    reacted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (message_id, user_id),
    FOREIGN KEY (message_id) REFERENCES public.message(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES public.user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS public.message_view (
    message_id UUID NOT NULL,
    user_id UUID NOT NULL,
    viewed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (message_id, user_id),
    FOREIGN KEY (message_id) REFERENCES public.message(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES public.user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS public.message_payload (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id UUID NOT NULL,
    file_path TEXT NOT NULL,
    FOREIGN KEY (message_id) REFERENCES public.message(id) ON DELETE CASCADE
);