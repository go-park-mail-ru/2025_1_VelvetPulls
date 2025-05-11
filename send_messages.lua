wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["Cookie"] = "token=982998aa-e3f5-4093-a459-5875aa128089"

messages_sent = 0
MAX_MESSAGES = 100000

function request()
  if messages_sent < MAX_MESSAGES then
    messages_sent = messages_sent + 1
    local body = string.format('{"message": "Message #%d from wrk"}', messages_sent)
    return wrk.format(nil, "/api/chat/1abc8a98-c288-4c19-a385-91c45758ae6a
5be21d8a-6f61-4bde-b847-35d7dda06647/messages", nil, body)
  else
    return nil
  end
end
