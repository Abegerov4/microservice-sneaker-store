CREATE TABLE IF NOT EXISTS chat_history (
    id         UUID PRIMARY KEY,
    session_id UUID        NOT NULL,
    user_id    TEXT        NOT NULL DEFAULT '',
    role       TEXT        NOT NULL CHECK (role IN ('user', 'assistant')),
    content    TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_chat_history_session_id ON chat_history (session_id, created_at DESC);
