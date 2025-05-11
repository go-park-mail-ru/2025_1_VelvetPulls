wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"

-- Get environment variables
local session_id = os.getenv("SESSION_ID")
local chat_id = os.getenv("CHAT_ID")

-- Validate environment variables
if not session_id then
    error("SESSION_ID is not set")
end

if not chat_id then
    error("CHAT_ID is not set")
end

-- Set cookie header PROPERLY
wrk.headers["Cookie"] = "token=" .. session_id

local counter = 0

function request()
    counter = counter + 1
    local path = "/api/chat/" .. chat_id .. "/messages"
    local body = string.format('{"message":"Test message %d"}', counter)
    return wrk.format("POST", path, wrk.headers, body)
end