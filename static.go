package main

import (
	"html/template"
	"net/http"
	"os"
	"log"
)


type PageData struct {
	PageTitle string
	OneApp []StatusStruct
}


func (pg *PageData) AddItem(item StatusStruct) []StatusStruct{
	pg.OneApp = append(pg.OneApp, item)
	return pg.OneApp 
}

func (a *App) DisplayHtml(w http.ResponseWriter, r *http.Request) {
	
	var tmp StatusStruct
	
	AppData := []StatusStruct{}

	data := PageData{
		PageTitle: "Hello There!",
		OneApp: AppData,
	}

	all_ids, err := a.GetAllID()
	if err != nil {
		log.Printf("Cannot get all IDs from DB: %v", err)
	}
	for _, i := range all_ids {
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
