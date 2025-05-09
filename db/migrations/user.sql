-- Замена "your_database" на имя вашей базы данных
GRANT CONNECT ON DATABASE mydb TO user_for_chat; -- Доступ на подключение к бд.
GRANT USAGE ON SCHEMA public TO user_for_chat;           -- Доступ к схеме (по умолчанию public)

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO user_for_chat; --Доступ на чтение и запись всех таблиц