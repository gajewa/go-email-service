package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"gitlab.jit.team/smurfs/email-service/controller"
	"gitlab.jit.team/smurfs/email-service/repository"
)

func main() {
	log.Println("Starting email service api")

	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln("Errors while connecting to DB : ", err)
		panic(err)
	}

	emailRepo := repository.NewEmailRepository(conn)
	emailController := controller.NewEmailController(emailRepo)
	emailController.Register(api)

	log.Fatalln(http.ListenAndServe(":8080", r))
}
