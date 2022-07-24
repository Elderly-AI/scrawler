package tasks

import (
	"github.com/jmoiron/sqlx"
)

type Database struct {
	conn *sqlx.DB
}

func New(conn *sqlx.DB) Database {
	return Database{
		conn,
	}
}

//func (d Database) CreateTask(ctx context.Context, task model.Task) {
//	query := `SELECT tag_id, external_id, tag_title, count(*) OVER() AS total FROM tags
//			      WHERE deleted_at IS NULL AND tag_title ~* CONCAT('\y', $3 :: TEXT)
//			      LIMIT $1 OFFSET $2`
//	err = d.conn.Select(&tags, query, limit, offset, search)
//	return
//}
