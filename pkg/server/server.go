package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/bsbsm/feeder/pkg/db"
	"github.com/gorilla/mux"
)

func SetSQLiteDatabase(db *db.SQLiteDatabase) {
	storage = db
}

func BlockingListen(port int) {
	if storage == nil {
		panic("Set SQLite database before start listen")
	}

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/scripts.js", jsHandler).Methods("GET")
	r.HandleFunc("/api/news", getNewsList).Methods("GET")
	r.HandleFunc("/api/news/{id}", getNewsByID).Methods("GET")
	r.HandleFunc("/api/feed", createFeedSource).Methods("PUT")

	r.Use(panicHandler, logMiddleware)

	a := ":" + strconv.Itoa(port)
	fmt.Printf("Listening at '%s'\n", a)

	http.ListenAndServe(a, r)
}
