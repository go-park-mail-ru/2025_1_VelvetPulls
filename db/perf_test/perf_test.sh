#!/bin/bash

HOST="localhost:8080"

# Случайный username
RAND_NUM=$RANDOM
USERNAME="wrkuser_$RAND_NUM"
PASSWORD="password123"

echo "Регистрируем пользователя: $USERNAME"

# Шаг 1: Регистрация пользователя
REGISTER_RESPONSE=$(curl -s -i -L -X POST http://$HOST/api/register \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\",\"name\":\"$USERNAME\"}")

echo "$REGISTER_RESPONSE" | grep -i Set-Cookie

# Извлекаем session_id из Set-Cookie
SESSION_ID=$(echo "$REGISTER_RESPONSE" | grep -i "Set-Cookie" | grep -o "token=[^;]*" | cut -d= -f2)

if [ -z "$SESSION_ID" ]; then
  echo "Ошибка: не удалось получить session_id"
  exit 1
fi

echo "Получен session_id: $SESSION_ID"

# Шаг 2: Создание чата (без multipart)
CHAT_PAYLOAD='{
  "type": "group",
  "title": "Project Team",
  "users": []
}'

CHAT_RESPONSE=$(curl -s -X POST http://$HOST/api/chat \
  -H "Cookie: token=$SESSION_ID" \
  -H "Content-Type: application/json" \
  -d "$CHAT_PAYLOAD")

echo "Ответ от сервера на создание чата:"
echo "$CHAT_RESPONSE"

CHAT_ID=$(echo "$CHAT_RESPONSE" | grep -o '"id":"[0-9a-f-]\+"' | head -n1 | cut -d: -f2 | tr -d '"')

if [ -z "$CHAT_ID" ]; then
  echo "Ошибка: не удалось получить chat_id"
  exit 1
fi

echo "Создан чат с chat_id: $CHAT_ID"

# Шаг 3: Запуск wrk
SESSION_ID=$SESSION_ID CHAT_ID=$CHAT_ID wrk -t2 -c10 -d50s -s ./db/perf_test/send_messages.lua -v http://$HOST
SESSION_ID=$SESSION_ID CHAT_ID=$CHAT_ID wrk -t2 -c10 -d50s -s ./db/perf_test/read_messages.lua -v http://$HOST
