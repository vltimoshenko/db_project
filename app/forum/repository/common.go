package repository

import (
	"fmt"
	"io/ioutil"

	"github.com/db_project/pkg/config"
	. "github.com/db_project/pkg/models"
	"github.com/db_project/pkg/sql_queries"

	"github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	DbConn *sqlx.DB
}

func (Rep *Repository) LoadSchemaSQL() error {
	if Rep.DbConn == nil {
		return pgx.ErrDeadConn
	}

	content, err := ioutil.ReadFile(config.DBSchema)
	if err != nil {
		return err
	}

	tx, err := Rep.DbConn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.Exec(string(content)); err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (Rep *Repository) Disconn() {
	if Rep.DbConn != nil {
		Rep.DbConn.Close()
		Rep.DbConn = nil
	}
}

func (r *Repository) ClearDB() error {
	_, err := r.DbConn.Exec(sql_queries.Clear)
	return err
}

func (r *Repository) GetStatus() (Status, error) {
	row := r.DbConn.QueryRowx(sql_queries.SelectDBStatus)

	var status Status
	err := row.StructScan(&status)
	if err != nil {
		fmt.Println(err)
		// return post, fmt.Errorf(messages.PostDoesNotExist)
	}
	// scanPost.Created = timetz.Format(time.RFC3339Nano)
	return status, nil
}
