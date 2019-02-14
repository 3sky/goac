package main

import (

	"net/http"
	"encoding/json"
)

type ErrorMsg struct {
	CODE int `json:"ERROR_CODE"`
	MSG string `json:"MSG"`
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

		username, password, authOK := r.BasicAuth()
		
		if authOK == false {
			e := &ErrorMsg{CODE: 401, MSG: "Not authorized"}
			http.Error(w, "", 401)
			json.NewEncoder(w).Encode(e)
			return
		}
	
		if username != u.username || password != u.password {
			e := &ErrorMsg{CODE: 401, MSG: "Not authorized"}
			http.Error(w, "", 401)
			json.NewEncoder(w).Encode(e)
			return
		}
	
		h.ServeHTTP(w, r)
	}
}