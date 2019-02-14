package main

import (

	"net/http"
	"os"
	"log"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"


)
type User struct {
	username  string
	password  string
}

type App struct {
    Router *mux.Router
    DB     *gorm.DB
}

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

func (a *App) Run(addr string) {
	loggedRouter := handlers.LoggingHandler(os.Stdout, a.Router)
	http.ListenAndServe(addr, loggedRouter)
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
	a.Router.HandleFunc("/", a.DisplayHtml).Methods("GET")

}


// func CheckErr(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }
