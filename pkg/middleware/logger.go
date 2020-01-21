package middleware

import (
	// "fmt"

	"fmt"
	"net/http"
	"time"
)

// type Logger struct {
// 	StdLogger *log.Logger
// 	f         *os.File
// }

func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		fmt.Printf("LOG START [%s] %s\n",
			r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		fmt.Printf("LOG END [%s] %s %s\n",
			r.Method, r.URL.Path, time.Since(start))
	})
}

// func NewLogger() Logger {
// 	loger := Logger{}

// 	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
// 	if err != nil {
// 		log.Fatalf("error opening file: %v", err)
// 	}
// 	loger.f = f

// 	log.SetOutput(loger.f)
// 	return loger
// }
