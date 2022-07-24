-- +goose Up
-- +goose StatementBegin

CREATE FUNCTION task_status_started()
    RETURNS SMALLINT
    IMMUTABLE
    LANGUAGE SQL AS
'SELECT 0 :: SMALLINT';

CREATE FUNCTION task_status_failed()
    RETURNS SMALLINT
    IMMUTABLE
    LANGUAGE SQL AS
'SELECT 1 :: SMALLINT';

CREATE FUNCTION task_status_finished()
    RETURNS SMALLINT
    IMMUTABLE
    LANGUAGE SQL AS
'SELECT 2 :: SMALLINT';

CREATE TABLE IF NOT EXISTS tasks (
    task_id    TEXT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS tasks_logs (
    task_log_id SERIAL PRIMARY KEY,
    task_id     TEXT NOT NULL,
    status      SMALLINT NOT NULL DEFAULT task_status_started(),
    error       TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMP NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMP DEFAULT NULL,
    CONSTRAINT  fk_task_id FOREIGN KEY(task_id) REFERENCES tasks(task_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE tasks;
DROP TABLE tasks_logs;
DROP FUNCTION task_status_started;
DROP FUNCTION task_status_failed;
DROP FUNCTION task_status_finished;

-- +goose StatementEnd