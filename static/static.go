package static

import (
	"p2go/db"
	"html/template"
	"net/http"
	"os"
	"strings"
)


type PageData struct {
	PageTitle string
	OneApp []db.StatusStruct
}


func (pg *PageData) AddItem(item db.StatusStruct) []db.StatusStruct{
	pg.OneApp = append(pg.OneApp, item)
	return pg.OneApp 
}

func DisplayHtml(w http.ResponseWriter, r *http.Request) {
	
	var tmp db.StatusStruct
	var pwd string 
	
	AppData := []db.StatusStruct{}

	data := PageData{
		PageTitle: "Hello There!",
		OneApp: AppData,
	}

	all_ids := db.GetAllID()
	for _, i := range all_ids {
		tmp = db.SelectFromDBWhereID(int64(i))
		data.AddItem(tmp)
	}

	p, err := os.Getwd()
	db.CheckErr(err)

	if strings.Contains(p, "static") {
		pwd = p + "/hello.html"
	} else {
		pwd = p + "/static/hello.html"
	}

	tmpl := template.Must(template.ParseFiles(pwd))
	tmpl.Execute(w, data)
}
