package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"

	// "github.com/valyala/fasthttp"
	// "github.com/valyala/fasthttp/fasthttpadaptor"

	"github.com/db_project/app/forum/repository"
	"github.com/db_project/app/forum/service"
	"github.com/db_project/app/server/delivery"

	"github.com/db_project/pkg/config"
	"github.com/db_project/pkg/middleware"
)

func NewRouter() (*mux.Router, error) {
	router := mux.NewRouter()

	DBConn, err := OpenSqlxViaPgxConnPool(config.DBPath)
	if err != nil {
		return nil, err
	}

	uS := service.Service{
		Repository: &repository.Repository{
			DbConn: DBConn,
		},
	}

	if err != nil {
		return nil, err
	}

	h := delivery.Handler{
		Service: uS,
	}

	router = router.PathPrefix("/api").Subrouter()
	router.Use(middleware.AccessLogMiddleware)

	router.HandleFunc("/forum/create", h.CreateForum).Methods(http.MethodPost)
	router.HandleFunc("/forum/{slug}/create", h.CreateThread).Methods(http.MethodPost)

	router.HandleFunc("/forum/{slug}/details", h.GetForum).Methods(http.MethodGet)
	router.HandleFunc("/forum/{slug}/threads", h.GetThreads).Methods(http.MethodGet)
	router.HandleFunc("/forum/{slug}/users", h.GetUsers).Methods(http.MethodGet)

	router.HandleFunc("/post/{id}/details", h.GetPost).Methods(http.MethodGet)
	router.HandleFunc("/post/{id}/details", h.ChangePost).Methods(http.MethodPost)

	router.HandleFunc("/user/{nickname}/create", h.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/user/{nickname}/profile", h.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/user/{nickname}/profile", h.ChangeUser).Methods(http.MethodPost)

	router.HandleFunc("/thread/{slug_or_id}/create", h.CreatePosts).Methods(http.MethodPost)
	router.HandleFunc("/thread/{slug_or_id}/details", h.GetThread).Methods(http.MethodGet)
	router.HandleFunc("/thread/{slug_or_id}/details", h.ChangeThread).Methods(http.MethodPost)

	router.HandleFunc("/service/clear", h.Clear).Methods(http.MethodPost)
	router.HandleFunc("/service/status", h.GetStatus).Methods(http.MethodGet)

	router.HandleFunc("/thread/{slug_or_id}/posts", h.GetPosts).Methods(http.MethodGet)
	router.HandleFunc("/thread/{slug_or_id}/vote", h.Vote).Methods(http.MethodPost)

	return router, nil
}

func RunServer() {
	router, err := NewRouter()
	if err != nil {
		log.Println(err.Error())
		log.Fatal("Failed to create router")
	}
	// log.Fatal(fasthttp.ListenAndServe(config.HostAddress, fasthttpadaptor.NewFastHTTPHandler(router)))
	log.Fatal(http.ListenAndServe(config.HostAddress, router))
}

// func OpenSqlxViaPgxConnPool(psqURI string) (*sqlx.DB, error) {
// 	conf, err := pgx.ParseURI(psqURI)
// 	if err != nil {
// 		return nil, err
// 	}

// 	connPool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
// 		ConnConfig: conf,
// 		// MaxConnections: config.MaxConn,
// 	})

// 	if err != nil {
// 		log.Println(err.Error())
// 		log.Fatal("Failed to create connections pool")
// 	}

// 	nativeDB := stdlib.OpenDBFromPool(connPool)

// 	fmt.Println("OpenSqlxViaPgxConnPool: the connection was created")
// 	// return sqlx.NewDb(nativeDB, "pgx"), nil
// 	return nativeDB
// }

func OpenSqlxViaPgxConnPool(psqURI string) (*sql.DB, error) {
	conf, err := pgx.ParseURI(psqURI)
	if err != nil {
		return nil, err
	}

	connPool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: conf,
		// MaxConnections: config.MaxConn,
	})

	if err != nil {
		log.Println(err.Error())
		log.Fatal("Failed to create connections pool")
	}

	nativeDB := stdlib.OpenDBFromPool(connPool)

	fmt.Println("OpenSqlxViaPgxConnPool: the connection was created")
	// return sqlx.NewDb(nativeDB, "pgx"), nil
	return nativeDB, nil
}
