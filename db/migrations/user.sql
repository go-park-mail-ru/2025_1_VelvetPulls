GRANT CONNECT ON DATABASE mydb TO user_for_chat; -- Доступ на подключение к бд.
GRANT USAGE ON SCHEMA public TO user_for_chat;   -- Доступ к схеме

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO user_for_chat; --Доступ на чтение и запись всех таблиц