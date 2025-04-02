package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/handler/auth"
	myctx "github.com/TechBowl-japan/go-stations/handler/context"
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

	userID := os.Getenv("BASIC_AUTH_USER_ID")
	if userID == "" {
		userID = "test"
	}
	pass := os.Getenv("BASIC_AUTH_PASSWORD")
	if pass == "" {
		pass = "test"
	}
	fmt.Println("userID =", userID, "\npassword =", pass)

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
	mux.Handle("/no-auth", handler.NewHealthzHandler())
	mux.Handle("/basic-auth", middleware.BasicAuth(handler.NewHealthzHandler(), *auth.NewConfigFromEnv()))
	mux.Handle("/slow", http.HandlerFunc(slowOperation))
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Println("ListenAndServe失敗:", err)
		}
	}()

	fmt.Println("起動済み")
	<-ctx.Done()
	fmt.Println("シャットダウン中...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("シャットダウン失敗: %v\n", err)
	} else {
		fmt.Println("シャットダウン成功")
	}
	return nil
}

// slowOperation
func slowOperation(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
	w.Write([]byte("slow operation done"))
}

// PanicHandler is a test handler that always panics.
type PanicHandler struct{}

func (p *PanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("panic test")
}

// OSCheckHandler is a test handler that checks the OS from the request context.
type OSCheckHandler struct{}

func (o *OSCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	os := myctx.OS(r.Context())
	if os == "" {
		http.Error(w, "os not found", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(os))
}

// curlの実行例
// curl.exe -u test:test http://localhost:8080/basic-auth
// envの設定例
// $env:BASIC_AUTH_USER_ID="hoge"; $env:BASIC_AUTH_PASSWORD="fuga"; go run .
