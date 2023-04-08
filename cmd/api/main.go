package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/twmb/franz-go/pkg/kgo"
	"log"
	"mborgnolo/companyservice/internal/data"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const version = "1.0.0"

// config holds the application configuration.
type config struct {
	port string
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

// application holds the dependencies for HTTP handlers.
type application struct {
	config  config
	logger  *log.Logger
	company interface {
		GetCompany(id uuid.UUID) (*data.Company, error)
		CreateCompany(company *data.Company) (uuid.UUID, error)
		DeleteCompany(id uuid.UUID) error
		UpdateCompany(company *data.Company) error
	}
	eventChan   chan data.EventRecord
	KafkaClient *kgo.Client
}

func main() {
	var cfg config
	flag.StringVar(&cfg.port, "port", os.Getenv("CMPSRV_PORT"), "API server port")
	flag.StringVar(&cfg.env, "env", os.Getenv("CMPSRV_ENV"), "Environment (development|testing|production)")
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
	// Initialize a new instance of application containing the dependencies.
	app := &application{
		config:      cfg,
		logger:      logger,
		company:     data.NewCompanyModel(db),
		eventChan:   make(chan data.EventRecord, 100),
		KafkaClient: initKafkaClient(),
	}
	// Initialize a new HTTP server.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.port),
		Handler:      app.routes(),
		ErrorLog:     nil,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)
	// Start a background goroutine that listens for SIGINT and SIGTERM signals
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		logger.Printf("shutting down server: %s:%s", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
		close(app.eventChan)
	}()

	logger.Println("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  cfg.env,
	})
	// Start a background goroutine that processes events.
	go app.processEvents()
	err = srv.ListenAndServe()

	logger.Fatal(err)
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

// processEvents is a background goroutine that processes events.
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
			message := fmt.Sprintf("company with id:[%s] %s at %s", event.ID, t, event.TimeStamp.Format(time.RFC3339))
			app.logger.Println(message)
			record := &kgo.Record{Topic: "companyservice", Value: []byte(message)}
			ctx := context.Background()
			app.KafkaClient.Produce(ctx, record, func(_ *kgo.Record, err error) {
				if err != nil {
					fmt.Printf("record had a produce error: %v\n", err)
				}
			})
		}
	}
}

// initKafkaClient initializes a new kafka client.
func initKafkaClient() *kgo.Client {
	seeds := []string{"kafka:9092"}
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.ConsumeTopics("companyservice"),
	)
	if err != nil {
		panic(err)
	}
	return cl
}
