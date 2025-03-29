## Таблица `user`

**Функциональные зависимости:**
`{id} → {avatar_path, first_name, last_name, username, phone, email, password, created_at, updated_at}`

Таблица находится в **1НФ**, так как все атрибуты атомарны, а строки уникальны.

Таблица находится во **2НФ**, так как нет частичной зависимости от составного ключа (первичный ключ - `id`), все атрибуты зависят от этого ключа.

Таблица находится в **3НФ**, так как все атрибуты нетранзитивно зависят от первичного ключа.

Таблица находится в **НФБК**, так как единственные детерминанты — это первичный ключ `id`.

## Таблица `chat`

**Функциональные зависимости:**
`{id} → {avatar_path, type, title, created_at, updated_at}`

Таблица находится в **1НФ**, так как все атрибуты атомарны, а строки уникальны.

Таблица находится во **2НФ**, так как нет частичной зависимости от составного ключа (первичный ключ - `id`).

Таблица находится в **3НФ**, так как все атрибуты нетранзитивно зависят от первичного ключа.

Таблица находится в **НФБК**, так как единственные детерминанты — это первичный ключ `id`.

## Таблица `user_chat`

**Функциональные зависимости:**
`{user_id, chat_id} → {user_role, joined_at}`

Таблица находится в **1НФ**, так как все атрибуты атомарны, а строки уникальны.

Таблица находится во **2НФ**, так как нет частичной зависимости от составного ключа (первичный ключ — `{user_id, chat_id}`). Все атрибуты зависят от этого ключа.

Таблица находится в **3НФ**, так как все атрибуты нетранзитивно зависят от первичного ключа.

Таблица находится в **НФБК**, так как единственные детерминанты — это составной первичный ключ `{user_id, chat_id}`.

## Таблица `contact`

**Функциональные зависимости:**
`{id} → {user_id, contact_id, created_at, updated_at}`

Таблица находится в **1НФ**, так как все атрибуты атомарны, а строки уникальны.

Таблица находится во **2НФ**, так как нет частичной зависимости от составного ключа (первичный ключ — `id`).

Таблица находится в **3НФ**, так как все атрибуты нетранзитивно зависят от первичного ключа.

Таблица находится в **НФБК**, так как единственные детерминанты — это первичный ключ `id`.

## Таблица `message`

**Функциональные зависимости:**
`{id} → {parent_message_id, chat_id, user_id, body, sent_at, is_redacted}`

Таблица находится в **1НФ**, так как все атрибуты атомарны, а строки уникальны.

Таблица находится во **2НФ**, так как нет частичной зависимости от составного ключа (первичный ключ — `id`).

Таблица находится в **3НФ**, так как все атрибуты нетранзитивно зависят от первичного ключа.

Таблица находится в **НФБК**, так как единственные детерминанты — это первичный ключ `id`.

## Таблица `message_reaction`

**Функциональные зависимости:**
`{id} → {message_id, user_id, reaction, reacted_at}`  
`{message_id, user_id} → {reaction, reacted_at}`

Таблица находится в **1НФ**, так как все атрибуты атомарны, а строки уникальны.

Таблица находится во **2НФ**, так как нет частичной зависимости от составного ключа (первичный ключ — `id`), и вторая зависимость `{message_id, user_id} → {reaction, reacted_at}` является потенциальным ключом.

Таблица находится в **3НФ**, так как все атрибуты нетранзитивно зависят от первичного ключа.

Таблица находится в **НФБК**, так как единственные детерминанты — это `{message_id, user_id}` (потенциальный ключ).

## Таблица `message_view`

**Функциональные зависимости:**
`{message_id, user_id} → {viewed_at}`

Таблица находится в **1НФ**, так как все атрибуты атомарны, а строки уникальны.

Таблица находится во **2НФ**, так как нет частичной зависимости от составного ключа (первичный ключ — `{message_id, user_id}`).

Таблица находится в **3НФ**, так как все атрибуты нетранзитивно зависят от первичного ключа.

Таблица находится в **НФБК**, так как единственные детерминанты — это составной первичный ключ `{message_id, user_id}`.

## Таблица `message_payload`

**Функциональные зависимости:**
`{id} → {message_id, file_path}`

Таблица находится в **1НФ**, так как все атрибуты атомарны, а строки уникальны.

Таблица находится во **2НФ**, так как нет частичной зависимости от составного ключа (первичный ключ — `id`).

Таблица находится в **3НФ**, так как все атрибуты нетранзитивно зависят от первичного ключа.

Таблица находится в **НФБК**, так как единственные детерминанты — это первичный ключ `id`.

```mermaid
erDiagram
    USER ||--o{ CONTACT : "has contacts"
    CHAT ||--o{ MESSAGE : "contains"
    USER ||--o{ MESSAGE : "writes"
    CHAT ||--o{ USER_CHAT : "has members"
    USER ||--o{ USER_CHAT : "participates"
    MESSAGE ||--o{ MESSAGE_PAYLOAD : "has attachment"
    MESSAGE ||--o{ MESSAGE_REACTION : "receives reactions"
    USER ||--o{ MESSAGE_REACTION : "reacts"
    MESSAGE ||--o{ MESSAGE_VIEW : "is viewed"
    USER ||--o{ MESSAGE_VIEW : "views"
    MESSAGE ||--o{ MESSAGE : "replies to"

    USER {
        uuid id PK
        text avatar_path
        text first_name
        text last_name
        text username
        text phone
        text email
        text password
        timestamp created_at
        timestamp updated_at
    }

    CHAT {
        uuid id PK
        text avatar_path
        chat_type type
        text title
        timestamp created_at
        timestamp updated_at
    }

    USER_CHAT {
        uuid user_id PK,FK
        uuid chat_id PK,FK
        chat_type user_role
        timestamp joined_at
    }

    CONTACT {
        uuid id PK
        uuid user_id FK
        uuid contact_id FK
        timestamp created_at
        timestamp updated_at
    }

    MESSAGE {
        uuid id PK
        uuid parent_message_id FK
        uuid chat_id FK
        uuid user_id FK
        text body
        timestamp sent_at
        bool is_redacted
    }

    MESSAGE_REACTION {
        uuid id PK
        uuid message_id FK
        uuid user_id FK
        text reaction
        timestamp reacted_at
    }

    MESSAGE_VIEW {
        uuid message_id PK,FK
        uuid user_id PK,FK
        timestamp viewed_at
    }

    MESSAGE_PAYLOAD {
        uuid id PK
        uuid message_id FK
        text file_path
    }
```