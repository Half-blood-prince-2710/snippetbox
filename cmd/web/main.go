package main

import (
	"database/sql"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/half-blood-prince-2710/snippetbox/internal/models"
)

type application struct {
	logger   *slog.Logger
	snippets *models.SnippetModel
	sessionManager *scs.SessionManager
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		// Use debug.Stack() to get the stack trace. This returns a byte slice, which
		// we need to convert to a string so that it's readable in the log entry.
		trace = string(debug.Stack())
	)
	// Include the trace in the log entry.
	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func main() {

	addr := flag.String("addr", ":5001", "HTTp network address")
	dsn := flag.String("dsn", "web:Manish@62229@/snippetbox?parseTime=true", "MySql data source name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()


	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	app := &application{
		logger:   logger,
		snippets: &models.SnippetModel{DB: db},
		sessionManager : sessionManager,
	}
	srv := &http.Server{
		Addr:		 *addr,
		Handler: app.routes(),
		ErrorLog: slog.NewLogLogger(logger.Handler(),slog.LevelError),	
	}
	logger.Info("Starting server on", "addr", *addr,"session",sessionManager.Cookie)

	log.Print("server is running on a port i wont tell you")

	err = srv.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
