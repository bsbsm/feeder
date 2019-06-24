package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/bsbsm/feeder/pkg/db"
	"github.com/gorilla/mux"
)

const maxCountParamValue = 100

var storage *db.SQLiteDatabase

func getNewsList(w http.ResponseWriter, r *http.Request) {
	offsetStr := r.URL.Query().Get("off")
	countStr := r.URL.Query().Get("c")
	titleSearchStr := r.URL.Query().Get("t")

	var offset, count int
	var err error

	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			panic(err)
		}
	}

	if countStr != "" {
		count, err = strconv.Atoi(countStr)
		if err != nil {
			panic(err)
		}
	}

	if count <= 0 {
		count = 10
	} else if count > maxCountParamValue {
		count = maxCountParamValue
	}

	var result []*db.News

	if titleSearchStr == "" {
		result, err = storage.GetNews(offset, count)
	} else {
		result, err = storage.GetNewsWithTitle(titleSearchStr, offset, count)
	}

	if err != nil {
		panic(err)
	}

	rsp, err := json.Marshal(result)

	if err != nil {
		panic(err)
	}

	fmt.Printf("\tJSON response:\n%s\n", string(rsp))

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(rsp); err != nil {
		panic(err)
	}
}

func getNewsByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		panic(err)
	}

	d, err := storage.GetNewsDetail(id)

	if err != nil {
		panic(err)
	}

	rsp, err := json.Marshal(d)

	if err != nil {
		panic(err)
	}

	fmt.Printf("\tJSON response:\n%s\n", string(rsp))

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(rsp); err != nil {
		panic(err)
	}
}

func createFeedSource(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("u")
	rule := r.URL.Query().Get("r")

	if err := storage.CreateFeedSource(url, rule); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../web/index.html")
}

func jsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	http.ServeFile(w, r, "../web/scripts.js")
}

func panicHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				const size = 64 << 10
				var b = make([]byte, size)
				b = b[:runtime.Stack(b, false)]
				fmt.Println(fmt.Sprintf("Client %s panic while serve request %s: %s", r.RemoteAddr, r.URL, err))
				panic(err)
			}
		}()
		h.ServeHTTP(w, r)
	})
}

func logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			fmt.Printf("Serving '%s' in %v\n", r.RequestURI, time.Since(start))
		}()
		h.ServeHTTP(w, r)
	})
}
