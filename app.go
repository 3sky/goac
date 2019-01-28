package main

import (
	"p2go/db"
	"p2go/api"
	st "p2go/static"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"encoding/base64"

)
type User struct {
	username  string
	password  string
}

func main() {

	router := mux.NewRouter()
	
	db.MakeMIgration()
	
	p := &User{
		username: "3sky",
		password: "test",
	
	}

	router.HandleFunc("/hello", use(api.SayHello, p.basicAuth)).Methods("GET")
	router.HandleFunc("/api/app/{id}", use(api.DisplaAppByID, p.basicAuth)).Methods("GET")
	router.HandleFunc("/api/app/new", use(api.AddNewApp, p.basicAuth)).Methods("POST")
	router.HandleFunc("/api/app/{id}", use(api.UpdateData, p.basicAuth)).Methods("PUT")
	router.HandleFunc("/api/app/{id}", use(api.DeleteData, p.basicAuth)).Methods("DELETE")
	router.HandleFunc("/api/apps", use(api.DisplayAllApp, p.basicAuth)).Methods("GET")
	router.HandleFunc("/", st.DisplayHtml).Methods("GET")
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	http.ListenAndServe(":5000", loggedRouter)
}


func use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}

	return h
}

func (u *User) basicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		if pair[0] != u.username || pair[1] != u.password {
			http.Error(w, "Not authorized", 401)
			return
		}

		h.ServeHTTP(w, r)
	}
}