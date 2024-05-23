package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bartoszjasak/api"
	sqlc "github.com/bartoszjasak/db/sqlc/generated"
	"github.com/bartoszjasak/service"
	"github.com/kelseyhightower/envconfig"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

type ConfigDB struct {
	Username string `envconfig:"DB_USERNAME"`
	Password string `envconfig:"DB_PASSWORD"`
	Port     string `envconfig:"DB_PORT"`
	Name     string `envconfig:"DB_NAME"`
	Host     string `envconfig:"DB_HOST"`
}

func main() {
	time.Sleep(time.Second * 1)

	var dbConf ConfigDB
	if err := envconfig.Process("MAILING", &dbConf); err != nil {
		log.Fatalln("Error loading environment variables", err)
	}

	dbConn, close := connectDB(dbConf)
	defer close()

	mailingService := service.New(sqlc.New(dbConn))
	mailingAPI := api.New(mailingService)

	r := chi.NewRouter()
	r.Use(middleware.Timeout(60 * time.Second))

	srv := http.Server{}
	srv.Addr = ":8080"
	srv.Handler = mailingAPI.NewHandler(r)

	wg := sync.WaitGroup{}
	shutdown, wait := gracefulShutdown(&srv, &wg)

	mailingService.StartCleanupJob(&wg, shutdown)

	log.Println("Starting http server")
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("Listen and serve returned error: %v", err)
	}

	<-wait
}

func connectDB(c ConfigDB) (*sql.DB, func() error) {
	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Name)

	dbConn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Printf("Postgres connect error : (%v)", err)
	}

	return dbConn, dbConn.Close
}

func gracefulShutdown(srv *http.Server, wg *sync.WaitGroup) (<-chan struct{}, <-chan struct{}) {
	shutdown := make(chan struct{})
	wait := make(chan struct{})
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		doneChan := make(chan struct{})
		go func() {
			wg.Wait()
			close(doneChan)
		}()

		close(shutdown)

		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown returned an err: %v\n", err)
		}

		select {
		case <-shutdownCtx.Done():
			log.Println("Graceful shutdown finnished with timeout")
		case <-doneChan:
			log.Println("Graceful shutdown finnished succesfully")
		}

		close(wait)
	}()

	return shutdown, wait
}
