package main

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// User for Basic Auth
type User struct {
	username string
	password string
}

// App init structure
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

//Initialize the app
func (a *App) Initialize(dbname string) {

	var err error
	a.DB, err = gorm.Open("sqlite3", dbname)
	CheckErr(err)

	a.MakeMigration()
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// Run the app on port
func (a *App) Run(addr string) {
	loggedRouter := handlers.LoggingHandler(os.Stdout, a.Router)
	http.ListenAndServe(addr, loggedRouter)
}

func (a *App) initializeRoutes() {

	p := &User{
		username: "3sky",
		password: "test",
	}

	a.Router.HandleFunc("/hello", use(SayHello, p.basicAuth)).Methods("GET")
	a.Router.HandleFunc("/api/app/{id}", use(a.DisplaAppByID, p.basicAuth)).Methods("GET")
	a.Router.HandleFunc("/api/app/new", use(a.AddNewApp, p.basicAuth)).Methods("POST")
	a.Router.HandleFunc("/api/app/{id}", use(a.UpdateData, p.basicAuth)).Methods("PUT")
	a.Router.HandleFunc("/api/app/{id}", use(a.DeleteData, p.basicAuth)).Methods("DELETE")
	a.Router.HandleFunc("/api/apps", use(a.DisplayAllApp, p.basicAuth)).Methods("GET")
	a.Router.HandleFunc("/", a.DisplayHTML).Methods("GET")

}

// CheckErr check errors
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
