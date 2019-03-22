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

// DisplayHTMLDev - Display HTML
func (a *App) DisplayHTMLDev(w http.ResponseWriter, r *http.Request) {

	var dev PageData

	dev, path, err := a.preperHTML("dev")

	if err != nil {
		log.Printf("cannot get dev aplication from db: %v", err)
	}

	tmpl := template.Must(template.ParseFiles(path + "/hello.html"))
	tmpl.Execute(w, dev)
}

// DisplayHTMLStg - Display HTML for stage
func (a *App) DisplayHTMLStg(w http.ResponseWriter, r *http.Request) {

	var stg PageData

	stg, path, err := a.preperHTML("stg")

	if err != nil {
		log.Printf("cannot get stg aplication from db: %v", err)
	}

	tmpl := template.Must(template.ParseFiles(path + "/hello.html"))
	tmpl.Execute(w, stg)
}

func (a *App) preperHTML(env string) (PageData, string, error) {

	var tmp StatusStruct

	appdata := []StatusStruct{}

	data := PageData{
		PageTitle: "Hello There!",
		OneApp:    appdata,
	}

	allids, err := a.GetAllID()
	if err != nil {
		log.Printf("cannot get all ids from db: %v", err)
	}
	for _, i := range allids {
		tmp, err = a.SelectFromDBWhereID(int64(i))
		if err != nil {
			log.Printf("cannot get row from db: %v", err)
		}
		if tmp.Environment == env {
			data.AddItem(tmp)
		}

	}

	p, err := os.Getwd()
	if err != nil {
		log.Printf("error while get path: %v", err)
	}

	return data, p, nil
}
