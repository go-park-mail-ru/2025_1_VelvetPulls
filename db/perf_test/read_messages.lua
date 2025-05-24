wrk.method = "GET"
wrk.headers["Content-Type"] = "application/json"

-- Получаем переменные окружения
local session_id = os.getenv("SESSION_ID")
local chat_id = os.getenv("CHAT_ID")

-- Проверяем переменные окружения
if not session_id then
    error("SESSION_ID is not set")
end

if not chat_id then
    error("CHAT_ID is not set")
end

-- Устанавливаем заголовок Cookie правильно
wrk.headers["Cookie"] = "token=" .. session_id

-- Функция для выполнения запроса
function request()
    local path = "/api/chat/" .. chat_id .. "/messages"
    return wrk.format("GET", path, wrk.headers)
end
