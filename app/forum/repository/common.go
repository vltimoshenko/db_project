package repository

import (
	. "github.com/db_project/pkg/models"
	"github.com/db_project/pkg/sql_queries"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	DbConn *sqlx.DB
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
