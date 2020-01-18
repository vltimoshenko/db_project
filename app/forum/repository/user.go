package repository

import (
	"fmt"

	"github.com/db_project/pkg/messages"
	. "github.com/db_project/pkg/models"
	"github.com/db_project/pkg/sql_queries"
)

// func (r *Repository) GetUser(nickname string, email string) ([]User, error) {
// 	var users []User

// 	row := r.DbConn.QueryRowx(sql_queries.SelectUserByNickname, nickname)

// 	var user User
// 	err := row.StructScan(&user)
// 	if err != nil {
// 		fmt.Println(err)
// 		// return users, err //fmt.Errorf()
// 	} else {
// 		users = append(users, user)
// 	}

// 	// row = r.DbConn.QueryRowx(sql_queries.SelectUserByEmail, email)

// 	// err = row.StructScan(&user)
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// 	// return users, err //fmt.Errorf()
// 	// } else {
// 	// 	users = append(users, user)
// 	// }

// 	return users, nil
// }

func (r *Repository) GetUserByNickname(nickname string) (User, error) {

	row := r.DbConn.QueryRowx(sql_queries.SelectUserByNickname, nickname)

	var user User
	err := row.StructScan(&user)
	if err != nil {
		fmt.Println(err)
		return user, err
	}

	return user, nil
}

func (r *Repository) GetUserByEmail(email string) (User, error) {

	row := r.DbConn.QueryRowx(sql_queries.SelectUserByEmail, email)

	var user User
	err := row.StructScan(&user)
	if err != nil {
		fmt.Println(err)
		return user, err
	}

	return user, nil
}

func (r *Repository) CreateUser(user NewUser, nickname string) error {
	var id int
	row, err := r.DbConn.Exec(sql_queries.InsertUser, user.About,
		user.Email, user.Fullname, nickname)

	fmt.Println(row)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf(messages.UserAlreadyExists)
	}

	fmt.Println(id)
	return nil
}
