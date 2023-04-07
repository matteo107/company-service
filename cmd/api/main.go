package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"log"
	"mborgnolo/companyservice/internal/data"
	"net/http"
	"os"
	"time"
)

const version = "0.0.1"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	jwt struct {
		secret string
	}
}

type application struct {
	config  config
	logger  *log.Logger
	company interface {
		GetCompany(id uuid.UUID) (*data.Company, error)
		CreateCompany(company *data.Company) (uuid.UUID, error)
		DeleteCompany(id uuid.UUID) error
		UpdateCompany(company *data.Company) error
	}
	eventChan chan data.EventRecord
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DB_DSN"), "PostgreSQL DSN")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	flag.StringVar(&cfg.jwt.secret, "jwt-secret", os.Getenv("JWT_SECRET"), "JWT secret")
	flag.Parse()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Printf("database connection pool established")
	app := &application{
		config:    cfg,
		logger:    logger,
		company:   data.NewCompanyModel(db),
		eventChan: make(chan data.EventRecord, 100),
	}
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		ErrorLog:     nil,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	logger.Println("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  cfg.env,
	})
	go app.processEvents()
	err = srv.ListenAndServe()
	//err = srv.ListenAndServeTLS("./certs/cert.pem", "./certs/key.pem")

	logger.Fatal(err, nil)
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
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

func (app *application) processEvents() {
	var t string
	for {
		select {
		case event := <-app.eventChan:
			switch event.Type {
			case data.CompanyCreated:
				t = "created"
			case data.CompanyDeleted:
				t = "deleted"
			case data.CompanyUpdated:
				t = "updated"
			}
			app.logger.Printf("company with id:[%s] %s at %s", event.ID, t, event.TimeStamp.Format(time.RFC3339))
		}
	}
}
