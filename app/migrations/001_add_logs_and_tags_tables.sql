-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tags (
    tag_id        SERIAL PRIMARY KEY,
    external_id   BIGINT NOT NULL,
    tag_title     TEXT NOT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT now(),
    updated_at    TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMP DEFAULT NULL,
    CONSTRAINT    unique_external_id_constraint UNIQUE (external_id, tag_title)
);

CREATE TABLE IF NOT EXISTS lesson_logs (
    lesson_log_id   SERIAL PRIMARY KEY,
    external_id     BIGINT NOT NULL,
    lessons_count   DECIMAL NOT NULL DEFAULT 0,
    created_at      TIMESTAMP NOT NULL DEFAULT now(),
    updated_at      TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMP DEFAULT NULL
);

CREATE UNIQUE INDEX unique_lesson_log_index ON lesson_logs(external_id, date_trunc('day', created_at));

CREATE TABLE IF NOT EXISTS lesson_log_tag (
    id            SERIAL PRIMARY KEY,
    tag_id        BIGINT NOT NULL,
    lesson_log_id BIGINT NOT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT now(),
    updated_at    TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMP DEFAULT NULL,
    CONSTRAINT    fk_tag_id FOREIGN KEY(tag_id) REFERENCES tags(tag_id),
    CONSTRAINT    fk_lesson_log_id FOREIGN KEY(lesson_log_id) REFERENCES lesson_logs(lesson_log_id)
);

CREATE UNIQUE INDEX unique_lesson_log_tag_index ON lesson_log_tag(tag_id, lesson_log_id, date_trunc('day', created_at));

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE tags;
DROP TABLE lesson_log;
DROP TABLE lesson_log_tag;

-- +goose StatementEnd