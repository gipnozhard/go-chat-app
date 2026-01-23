-- +goose Up
-- +goose StatementBegin

-- Этот раздел выполняется при применении миграции (goose up)
-- Здесь описываем все изменения, которые нужно внести в БД

-- Создаем таблицу "чаты" для хранения информации о чатах
CREATE TABLE chats (
                       id SERIAL PRIMARY KEY,         -- Уникальный идентификатор, автоматически увеличивается
                       title VARCHAR(200) NOT NULL,   -- Заголовок чата, обязательное поле, максимум 200 символов
                       created_at TIMESTAMP DEFAULT NOW() -- Дата и время создания, по умолчанию текущее время
);

-- Создаем таблицу "сообщения" для хранения текста сообщений
CREATE TABLE messages (
                          id SERIAL PRIMARY KEY,         -- Уникальный идентификатор сообщения
                          chat_id INTEGER NOT NULL REFERENCES chats(id) ON DELETE CASCADE, -- Ссылка на чат
    -- ^ chat_id - ID чата, к которому относится сообщение
    -- ^ REFERENCES chats(id) - внешний ключ на таблицу chats
    -- ^ ON DELETE CASCADE - при удалении чата автоматически удаляются все его сообщения
                          text TEXT NOT NULL,            -- Текст сообщения, обязательное поле
                          created_at TIMESTAMP DEFAULT NOW() -- Дата и время создания сообщения
);

-- Создаем индекс для ускорения поиска сообщений по ID чата
-- Без индекса поиск сообщений конкретного чата был бы медленным
CREATE INDEX ON messages(chat_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Этот раздел выполняется при откате миграции (goose down)
-- Здесь описываем как отменить изменения из раздела Up

-- Удаляем таблицу сообщений (в обратном порядке создания)
-- Сначала messages, так как она зависит от chats (внешний ключ)
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS chats;
-- +goose StatementEnd