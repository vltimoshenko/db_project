package repository

import (
	"fmt"

	. "github.com/db_project/pkg/models"
	"github.com/db_project/pkg/sql_queries"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	DbConn *sqlx.DB
	// DbConn *sql.DB
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
	return status, nil
}
