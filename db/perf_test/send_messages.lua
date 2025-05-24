wrk.method = "POST"

-- Получаем переменные окружения
local session_id = os.getenv("SESSION_ID")
local chat_id = os.getenv("CHAT_ID")

-- Проверка переменных
if not session_id then
    error("SESSION_ID is not set")
end

if not chat_id then
    error("CHAT_ID is not set")
end

-- Устанавливаем заголовок Cookie
wrk.headers["Cookie"] = "token=" .. session_id

-- Уникальная граница для multipart
local boundary = "------------------------abcd1234"
wrk.headers["Content-Type"] = "multipart/form-data; boundary=" .. boundary

-- Счётчик сообщений
local counter = 0

-- Печатаем один раз
print("SESSION_ID: " .. session_id)
print("CHAT_ID: " .. chat_id)

function request()
    counter = counter + 1
    local path = "/api/chat/" .. chat_id .. "/messages"
    
    local json_value = string.format('{"message":"Test message %d"}', counter)
    local body = "--" .. boundary .. "\r\n" ..
                 'Content-Disposition: form-data; name="text"' .. "\r\n\r\n" ..
                 json_value .. "\r\n" ..
                 "--" .. boundary .. "--\r\n"

    return wrk.format("POST", path, wrk.headers, body)
end