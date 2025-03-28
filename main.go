package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler/context"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/handler/router"
)

func main() {
	err := realMain()
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}
}

func realMain() error {
	// config values
	const (
		defaultPort   = ":8080"
		defaultDBPath = ".sqlite3/todo.db"
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	// set time zone
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	// set up sqlite3
	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		return err
	}
	defer todoDB.Close()

	// NOTE: 新しいエンドポイントの登録はrouter.NewRouterの内部で行うようにする
	mux := router.NewRouter(todoDB)
	// mux.Handle("/do-panic", &PanicHandler{})
	// mux.Handle("/do-panic2", middleware.Recovery(&PanicHandler{}))
	// mux.Handle("/os", middleware.WithOS(&OSCheckHandler{}))
	mux.Handle("/log", middleware.Log(&OSCheckHandler{}))
	mux.Handle("/log2", middleware.WithOS(middleware.Log(&OSCheckHandler{})))

	if err := http.ListenAndServe(port, mux); err != nil {
		return err
	}

	return nil
}

// PanicHandler is a test handler that always panics.
type PanicHandler struct{}

func (p *PanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("panic test")
}

// OSCheckHandler is a test handler that checks the OS from the request context.
type OSCheckHandler struct{}

func (o *OSCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	os := context.OS(r.Context())
	if os == "" {
		http.Error(w, "os not found", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(os))
}
