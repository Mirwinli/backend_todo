CREATE TABLE IF NOT EXISTS tasks (
    id          BIGSERIAL PRIMARY KEY,           -- Унікальний ID задачі
    user_id     BIGINT NOT NULL,                -- ID користувача з твого Auth-сервісу
    title       VARCHAR(255) NOT NULL,          -- Заголовок задачі
    description TEXT,  -- Опис (може бути NULL)
    duration    INTERVAL,
    is_done     BOOLEAN DEFAULT FALSE,          -- Статус виконання
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    done_at     TIMESTAMPTZ                    -- Буде NULL, поки задача не виконана
    );

CREATE INDEX idx_tasks_user_id ON tasks(user_id);