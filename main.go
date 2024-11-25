package main

import (
	"crud-basico/server"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/users", server.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/users", server.FindAllUsers).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", server.FindAllUsers).Methods(http.MethodGet)
	fmt.Println("The applicationa are running in the port: 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
