package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"
    "html/template" // New import
	_ "github.com/go-sql-driver/mysql" // New import
	"snippetbox.mcheng.net/internal/models"
)

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include the structured logger, but we'll
// add more to this as the build progresses.
type application struct{
    logger *slog.Logger
    snippets *models.SnippetModel
    templateCache map[string]*template.Template
}

func main() {



    // flag will be stored in the addr variable at runtime.
    addr := flag.String("addr", ":4000", "HTTP network address")
    dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
    flag.Parse()


    // Use the slog.New() function to initialize a new structured logger, which
    // writes to the standard out stream and uses the default settings.”
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        AddSource: true,
    }))

    db, err := openDB(*dsn)
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }

    defer db.Close()


    // Initialize a new template cache...
    templateCache, err := newTemplateCache()
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }

    // Initialize a new instance of our application struct, containing the
    // dependencies (for now, just the structured logger)”
    app := &application{
        logger: logger, 
        snippets: &models.SnippetModel{DB: db},
        templateCache: templateCache,
    }


    logger.Info("starting server", "addr",  *addr)

    err = http.ListenAndServe(*addr, app.routes())
    logger.Error(err.Error())
    os.Exit(1)

}

func openDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }

    err = db.Ping()
    if err != nil {
        db.Close()
        return nil, err
    }

    return db, nil
}



