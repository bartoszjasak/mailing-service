package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bartoszjasak/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

func main() {
	time.Sleep(time.Second * 1)

	r := chi.NewRouter()

	r.Use(middleware.Timeout(60 * time.Second))

	dbConn, close := ConnectDB()
	defer close()

	service := service.New(dbConn)

	r.Post("/api/messages", service.PostMessage)
	r.Post("/api/messages/send", service.SendMessages)
	r.Delete("/api/messages/{id}", service.DeleteMessage)

	log.Println("Starting http server")
	http.ListenAndServe(":8080", r)
}

func ConnectDB() (*sql.DB, func() error) {
	psqlInfo := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", "postgres", "postgres", "helm-postgresql", 5432, "postgres")
	dbConn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Printf("Postgres connect error : (%v)", err)
	}

	return dbConn, dbConn.Close
}
