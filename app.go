package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

//User BasicAuth User
type User struct {
	username string
	password string
}

//App - it's app
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

//Initialize - initialize App
func (a *App) Initialize(dbname string) {

	var err error
	a.DB, err = gorm.Open("sqlite3", dbname)
	if err != nil {
		log.Fatalf("Cannot open main DB %v", err)
	}

	err = a.MakeMigration()
	if err != nil {
		log.Fatalf("Cannot make schema migration: %v", err)
	}
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

//Run - run app
func (a *App) Run(addr string) {
	loggedRouter := handlers.LoggingHandler(os.Stdout, a.Router)
	/**
	To use https change http.ListenAndServe -> http.ListenAndServeTLS

	err := http.ListenAndServeTLS(addr, "server.pem", "server.key", loggedRouter)

	And create server.pem and server.key
	**/

	err := http.ListenAndServe(addr, loggedRouter)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
	defer a.DB.Close()
}

func (a *App) initializeRoutes() {

	p := &User{
		username: "3sky",
		password: "test",
	}

	a.Router.HandleFunc("/hello", use(SayHello, p.basicAuth)).Methods("GET")
	a.Router.HandleFunc("/api/app/{id}", use(a.DisplayAppByID, p.basicAuth)).Methods("GET")
	a.Router.HandleFunc("/api/app/new", use(a.AddNewApp, p.basicAuth)).Methods("POST")
	a.Router.HandleFunc("/api/app/{id}", use(a.UpdateData, p.basicAuth)).Methods("PUT")
	a.Router.HandleFunc("/api/app/{id}", use(a.DeleteData, p.basicAuth)).Methods("DELETE")
	a.Router.HandleFunc("/api/apps", use(a.DisplayAllApp, p.basicAuth)).Methods("GET")
	a.Router.HandleFunc("/dev", a.DisplayHTMLDev).Methods("GET")
	a.Router.HandleFunc("/stg", a.DisplayHTMLStg).Methods("GET")
}
