package main

import (
	"p2go/db"
	"p2go/api"
	st "p2go/static"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)


func main() {
	router := mux.NewRouter()
	db.MakeMIgration()
	router.HandleFunc("/hello", api.SayHello).Methods("GET")
	router.HandleFunc("/api/app/{id}", api.DisplaAppByID).Methods("GET")
	router.HandleFunc("/api/app/new", api.AddNewApp).Methods("POST")
	router.HandleFunc("/api/app/{id}", api.UpdateData).Methods("PUT")
	router.HandleFunc("/api/app/{id}", api.DeleteData).Methods("DELETE")
	router.HandleFunc("/api/apps", api.DisplayAllApp).Methods("GET")
	router.HandleFunc("/", st.DisplayHtml).Methods("GET")
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	http.ListenAndServe(":5000", loggedRouter)
}


