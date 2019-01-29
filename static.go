package main

import (
	"html/template"
	"net/http"
	"os"
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

	all_ids := a.GetAllID()
	for _, i := range all_ids {
		tmp = a.SelectFromDBWhereID(int64(i))
		data.AddItem(tmp)
	}

	p, err := os.Getwd()
	CheckErr(err)
	
	tmpl := template.Must(template.ParseFiles(p + "/hello.html"))
	tmpl.Execute(w, data)
}
