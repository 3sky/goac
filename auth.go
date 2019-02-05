package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
)

//ErrorMsg get  error message
type ErrorMsg struct {
	CODE int    `json:"ERROR_CODE"`
	MSG  string `json:"MSG"`
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
			e := &ErrorMsg{CODE: 401, MSG: "Not authorized"}
			http.Error(w, "", 401)
			json.NewEncoder(w).Encode(e)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			e := &ErrorMsg{CODE: 401, MSG: err.Error()}
			http.Error(w, "", 401)
			json.NewEncoder(w).Encode(e)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			e := &ErrorMsg{CODE: 401, MSG: "Not authorized"}
			http.Error(w, "", 401)
			json.NewEncoder(w).Encode(e)
			return
		}

		if pair[0] != u.username || pair[1] != u.password {
			e := &ErrorMsg{CODE: 401, MSG: "Not authorized"}
			http.Error(w, "", 401)
			json.NewEncoder(w).Encode(e)
			return
		}

		h.ServeHTTP(w, r)
	}
}
