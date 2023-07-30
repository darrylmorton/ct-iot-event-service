package app

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/darrylmorton/ct-iot-event-service/internal/data"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "0.0.1"

type config struct {
	port int
	env  string
	dsn  string
}

type application struct {
	config config
	logger *log.Logger
	models data.Models
}

type Event struct {
	Id          string    `json:"id,omitempty"`
	DeviceName  string    `json:"deviceName,omitempty"`
	Description string    `json:"description,omitempty"`
	Type        string    `json:"type,omitempty"`
	Event       string    `json:"event,omitempty"`
	Read        bool      `json:"read,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}

// Queue provides the ability to handle SQS messages.
type Queue struct {
	Client sqsiface.SQSAPI
	URL    string
}

type ValueUnit struct {
	Value int    `json:"value"`
	Unit  string `json:"unit"`
}

type Temperature struct {
	Value      int    `json:"value"`
	Unit       string `json:"unit"`
	Connection string `json:"connection"`
}

type Humidity struct {
	Value         int    `json:"value"`
	Unit          string `json:"unit"`
	Connection    string `json:"connection"`
	Precipitation bool   `json:"precipitation"`
}

type Payload struct {
	Cadence     ValueUnit   `json:"cadence"`
	Battery     ValueUnit   `json:"battery"`
	Temperature Temperature `json:"temperature"`
	Humidity    Humidity    `json:"humidity"`
}

// Message is a concrete representation of the SQS message
type Message struct {
	DeviceId         string  `json:"deviceId"`
	PayloadTimestamp string  `json:"payloadTimestamp"`
	Payload          Payload `json:"payload"`
}

func StartServer(q Queue) *http.Server {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev|stage|prod)")
	flag.StringVar(&cfg.dsn, "db-dsn", os.Getenv("EVENTS_DB_DSN"), "PostgreSQL DSN")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := sql.Open("postgres", cfg.dsn)
	if err != nil {
		logger.Fatal(err)
	}

	//defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Printf("database connection pool established")

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	go q.GetMessages(db, 20)

	addr := fmt.Sprintf(":%d", cfg.port)

	srv := &http.Server{
		Addr:         addr,
		Handler:      app.router(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	logger.Printf("starting %s server on %s", cfg.env, addr)

	return srv
}
