package main

import (
	"html/template"
	"net/http"
	"os"
)

// PageData struct with app info
type PageData struct {
	PageTitle string
	OneApp    []StatusStruct
}

//AddItem add app to list
func (pg *PageData) AddItem(item StatusStruct) []StatusStruct {
	pg.OneApp = append(pg.OneApp, item)
	return pg.OneApp
}

//DisplayHTML - Display Html
func (a *App) DisplayHTML(w http.ResponseWriter, r *http.Request) {

	var tmp StatusStruct

	AppData := []StatusStruct{}

	data := PageData{
		PageTitle: "Hello There!",
		OneApp:    AppData,
	}

	allIds := a.GetAllID()
	for _, i := range allIds {
		tmp = a.SelectFromDBWhereID(int64(i))
		data.AddItem(tmp)
	}

	p, err := os.Getwd()
	CheckErr(err)

	tmpl := template.Must(template.ParseFiles(p + "/hello.html"))
	tmpl.Execute(w, data)
}
