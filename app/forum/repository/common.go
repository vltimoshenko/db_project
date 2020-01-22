package repository

import (
	"database/sql"

	. "github.com/db_project/pkg/models"
	"github.com/db_project/pkg/sql_queries"
)

type Repository struct {
	DbConn *sql.DB
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
	row := r.DbConn.QueryRow(sql_queries.SelectDBStatus)
	var status Status
	_ = row.Scan(&status.Post, &status.Thread, &status.User, &status.Forum)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	return status, nil
}
