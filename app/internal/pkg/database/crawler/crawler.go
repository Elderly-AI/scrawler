package crawler

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"time"

	"github.com/Elderly-AI/scrawler/internal/pkg/model"
)

type Database struct {
	conn *sqlx.DB
}

func New(conn *sqlx.DB) Database {
	return Database{
		conn,
	}
}

func (d Database) AddTags(ctx context.Context, tags []model.Tag) (ids []uint64, err error) {
	externalIDs := make([]uint64, len(tags))
	tagTitles := make([]string, len(tags))
	for i, t := range tags {
		externalIDs[i] = t.ExternalID
		tagTitles[i] = t.Title
	}
	query := `INSERT INTO tags (external_id, tag_title) VALUES(UNNEST($1 :: BIGINT[]), UNNEST($2 :: TEXT[]))
			  ON CONFLICT (external_id, tag_title) DO UPDATE SET updated_at = now(), deleted_at = NULL
			  RETURNING tag_id`
	err = d.conn.Select(&ids, query, pq.Array(externalIDs), pq.Array(tagTitles))
	return
}

func (d Database) GetTagsWithPagination(ctx context.Context, search string, limit, offset uint64) (tags []model.TagDB, err error) {
	query := `SELECT tag_id, external_id, tag_title, count(*) OVER() AS total FROM tags
			      WHERE deleted_at IS NULL AND tag_title ~* CONCAT('\y', $3 :: TEXT)
			      LIMIT $1 OFFSET $2`
	err = d.conn.Select(&tags, query, limit, offset, search)
	return
}

func (d Database) AddLogsWithTags(ctx context.Context, logs []model.LessonLog) (err error) {
	tx, err := d.conn.BeginTx(ctx, &sql.TxOptions{})
	defer tx.Rollback()
	if err != nil {
		return err
	}

	logs, err = d.AddLogs(ctx, logs)
	if err != nil {
		return err
	}

	for _, l := range logs {
		var tags []model.TagDB
		tags, err = d.GetTagsByExternalIDs(ctx, l.Tags)
		if err != nil {
			return err
		}
		l.Tags = convertTags(tags)
		err = d.AddLogToTag(ctx, l)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (d Database) GetTagsByExternalIDs(ctx context.Context, tags []model.Tag) ([]model.TagDB, error) {
	tagsExternalIDs := make([]uint64, len(tags))
	for i, t := range tags {
		tagsExternalIDs[i] = t.ExternalID
	}
	var tg []model.TagDB
	query := `SELECT tag_id, external_id, tag_title FROM tags WHERE deleted_at IS NULL AND external_id = ANY ($1 :: BIGINT[])`
	err := d.conn.Select(&tg, query, pq.Array(tagsExternalIDs))
	return tg, err
}

func (d Database) GetLessonsCountByTags(ctx context.Context, tagsExternalIDs []uint64, from, to time.Time) (count float64, err error) {
	query := `SELECT SUM(ll.lessons_count) AS count FROM lesson_logs AS ll
         	  LEFT JOIN lesson_log_tag AS llt ON llt.lesson_log_id = ll.lesson_log_id
         	  LEFT JOIN tags AS t ON t.tag_id = llt.tag_id 
         	  WHERE t.deleted_at IS NULL AND llt.deleted_at IS NULL AND ll.deleted_at IS NULL AND
         	        t.external_id = ANY($1 :: BIGINT[]) AND ll.created_at > $2 AND ll.created_at < $3`
	err = d.conn.Get(&count, query, pq.Array(tagsExternalIDs), from, to)
	return
}

func (d Database) AddLogs(ctx context.Context, logs []model.LessonLog) ([]model.LessonLog, error) {
	externalIDs := make([]uint64, len(logs))
	lessonsCount := make([]float64, len(logs))
	for i, l := range logs {
		externalIDs[i] = l.ExternalID
		lessonsCount[i] = l.LessonsCount
	}

	var logsDB []model.LessonLogDB
	query := `INSERT INTO lesson_logs (external_id, lessons_count) 
			      VALUES (UNNEST($1 :: BIGINT[]), UNNEST($2 :: DECIMAL[])) 
			      ON CONFLICT (external_id, date_trunc('day', created_at)) DO NOTHING
			      RETURNING lesson_log_id, external_id, lessons_count`
	err := d.conn.Select(&logsDB, query, pq.Array(externalIDs), pq.Array(lessonsCount))
	if err != nil {
		return nil, err
	}

	logsDBMap := make(map[uint64]model.LessonLogDB)
	for i, l := range logsDB {
		logsDBMap[l.ExternalID] = logsDB[i]
	}

	var out []model.LessonLog
	for i, l := range logs {
		dbLog, ok := logsDBMap[l.ExternalID]
		if !ok {
			continue
		}
		tmpLog := logs[i]
		tmpLog.ID = dbLog.ID
		out = append(out, tmpLog)
	}
	return out, nil
}

func (d Database) AddLogToTag(ctx context.Context, log model.LessonLog) (err error) {
	tagIDs := make([]uint64, len(log.Tags))
	for i, t := range log.Tags {
		tagIDs[i] = t.ID
	}
	qq := `INSERT INTO lesson_log_tag (tag_id, lesson_log_id) 
		       VALUES (UNNEST($1 :: BIGINT[]), UNNEST(array_fill($2 :: INT, ARRAY[$3 :: INT]) :: BIGINT[]))
		       ON CONFLICT (tag_id, lesson_log_id, date_trunc('day', created_at)) DO NOTHING`
	_, err = d.conn.Exec(qq, pq.Array(tagIDs), log.ID, len(log.Tags))
	return err
}

func convertTags(tags []model.TagDB) []model.Tag {
	tg := make([]model.Tag, len(tags))
	for i, t := range tags {
		tg[i] = model.Tag{
			ID:         t.ID,
			ExternalID: t.ExternalID,
			Title:      t.Title,
		}
	}
	return tg
}
