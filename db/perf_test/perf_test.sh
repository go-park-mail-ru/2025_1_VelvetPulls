#!/bin/bash

HOST="localhost:8080"

# Случайный username и телефон
RAND_NUM=$RANDOM
USERNAME="wrkuser_$RAND_NUM"
PASSWORD="password123"
PHONE="7900$(printf "%07d" $((RANDOM % 10000000)))"

echo "Регистрируем пользователя: $USERNAME с телефоном: $PHONE"

# Шаг 1: Регистрация пользователя
REGISTER_RESPONSE=$(curl -s -i -X POST http://$HOST/api/register \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\",\"confirm_password\":\"$PASSWORD\",\"phone\":\"$PHONE\"}")

SESSION_ID=$(echo "$REGISTER_RESPONSE" | grep -Fi Set-Cookie | sed -n 's/.*token=\([^;]*\);.*/\1/p')

if [ -z "$SESSION_ID" ]; then
  echo "Ошибка: не удалось получить session_id"
  exit 1
fi

echo "Получен session_id: $SESSION_ID"

# Шаг 2: Создание чата
CHAT_RESPONSE=$(curl -s -X POST http://$HOST/api/chat \
  -H "Cookie: token=$SESSION_ID" \
  -H "Content-Type: multipart/form-data; boundary=------------------------abcd1234" \
  --data-binary $'--------------------------abcd1234\r\nContent-Disposition: form-data; name="chat_data"\r\n\r\n{"type":"group","title":"wrk_chat"}\r\n--------------------------abcd1234--\r\n')

echo "Ответ от сервера на создание чата:"
echo "$CHAT_RESPONSE"
CHAT_ID=$(echo "$CHAT_RESPONSE" | grep -o '"id":"[0-9a-f\-]\+"' | head -n1 | cut -d: -f2 | tr -d '"')

if [ -z "$CHAT_ID" ]; then
  echo "Ошибка: не удалось получить chat_id"
  exit 1
fi

echo "Создан чат с chat_id: $CHAT_ID"

echo "CHAT_ID=$CHAT_ID"
echo "SESSION_ID=$SESSION_ID"

# Шаг 3: Запуск wrk
SESSION_ID=$SESSION_ID CHAT_ID=$CHAT_ID wrk -t2 -c10 -d50s -s ./db/perf_test/send_messages.lua -v http://$HOST
SESSION_ID=$SESSION_ID CHAT_ID=$CHAT_ID wrk -t2 -c10 -d50s -s ./db/perf_test/read_messages.lua -v http://$HOST