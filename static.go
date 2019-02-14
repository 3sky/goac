package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

//PageData - Data for table in HTML
type PageData struct {
	PageTitle string
	OneApp    []StatusStruct
}

//AddItem - Add item to PageData struct
func (pg *PageData) AddItem(item StatusStruct) []StatusStruct {
	pg.OneApp = append(pg.OneApp, item)
	return pg.OneApp
}

// DisplayHTML - Display HTML
func (a *App) DisplayHTML(w http.ResponseWriter, r *http.Request) {

	var tmp StatusStruct

	AppData := []StatusStruct{}

	data := PageData{
		PageTitle: "Hello There!",
		OneApp:    AppData,
	}

	allIds, err := a.GetAllID()
	if err != nil {
		log.Printf("Cannot get all IDs from DB: %v", err)
	}
	for _, i := range allIds {
		tmp, err = a.SelectFromDBWhereID(int64(i))
		if err != nil {
			log.Printf("Cannot get row from DB: %v", err)
		}
		data.AddItem(tmp)
	}

	p, err := os.Getwd()
	if err != nil {
		log.Printf("Error while get path: %v", err)
	}

	tmpl := template.Must(template.ParseFiles(p + "/hello.html"))
	tmpl.Execute(w, data)
}
