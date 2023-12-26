package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/0xAckerMan/Lets-Go-Further/internal/data"
	_ "github.com/lib/pq"
)

const Version = "1.0.0"

type Config struct {
	env  string
	port int
	db   struct {
		dsn string
        maxOpenConns int
        maxIdleConns int
        maxIdleTime string
	}
}

type Application struct {
	logger *log.Logger
	config Config
    models data.Models
}

func main() {
	var cfg Config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "dev", "Running environment (dev|stag|prod)")
	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL dsn")
    flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
    flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
    flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max idle time")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	logger.Println("database connection pool established")

	app := &Application{
		config: cfg,
		logger: logger,
        models: data.NewModel(db),
	}

	addr := fmt.Sprintf(":%d", cfg.port)

	srv := &http.Server{
		Addr:         addr,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("Starting %s server on port %d ", cfg.env, cfg.port)

	srv.ListenAndServe()
}

func openDB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

    db.SetMaxOpenConns(cfg.db.maxOpenConns)
    db.SetMaxIdleConns(cfg.db.maxIdleConns)
    
    duration, err := time.ParseDuration(cfg.db.maxIdleTime)
    if err != nil{
        return nil, err
    }

    db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
