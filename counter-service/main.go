package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/galihrivanto/runner"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

var (
	service     = "counter-service"
	version     = "1.0.0"
	environment = "dev"
)

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	if name == "" {
		name = "Guest"
	}
	log.Printf("Received request for %s\n", name)
	w.Write([]byte(fmt.Sprintf("Hello, %s\n", name)))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	var (
		settings = new(Settings)
		ctx      = context.Background()
	)

	viper.SetConfigName("counter-service")
	viper.AddConfigPath("/etc/counter-service")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err := viper.Unmarshal(&settings); err != nil {
		panic(fmt.Errorf("failed to scan settings: %w", err))
	}

	// setup and teardown telemetry
	if err := setupTracerProvider(settings); err != nil {
		panic(err)
	}

	defer teardownTracerProvider(ctx)

	// open close for redis
	if err := OpenRedis(settings); err != nil {
		panic(err)
	}
	defer CloseRedis()

	// Create Server and Route Handlers
	r := mux.NewRouter()

	if settings.Telemetry.Enabled {
		mw := otelmux.Middleware(service)
		r.Use(mw)
	}

	r.HandleFunc("/", handler)
	r.HandleFunc("/health", healthHandler)
	r.HandleFunc("/readiness", readinessHandler)

	r.HandleFunc("/api/v1/counter/add", countAddHandler).Methods("POST")
	r.HandleFunc("/api/v1/counter/dec", countDecHandler).Methods("POST")
	r.HandleFunc("/api/v1/counter", countGetHandler).Methods("GET")

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())

	runner.
		Run(ctx, func(ctx context.Context) error {
			log.Println("Starting Server at :8080")
			if err := srv.ListenAndServe(); err != nil {
				return err
			}

			return nil
		}).
		Handle(func(sig os.Signal) {
			if sig == syscall.SIGHUP {
				return
			}

			log.Println("Shutting down...")
			srv.Shutdown(ctx)

			cancel()
		})
}
