package main

import (
	"html/template"
	"log"
	"net/http"
)

//PageData - Data for table in HTML
type PageData struct {
	PageTitle string
	OneApp    []AppStatusStruct
}

//AddItem - Add item to PageData struct
func (pg *PageData) AddItem(item AppStatusStruct) []AppStatusStruct {
	pg.OneApp = append(pg.OneApp, item)
	return pg.OneApp
}

const htm = `
<!doctype html>
<html>
<body>
<style>

table {
  font-family: arial, sans-serif;
  border-collapse: collapse;
  width: 100%;
}

td, th {
  border: 1px solid #dddddd;
  text-align: left;
  padding: 8px;
}

tr:nth-child(even) {
  background-color: #dddddd;
}
</style>
<h1 align="center"> {{.PageTitle}} </h1>
<table style="width:100%">
    <tr>
      <th>ID</th>
      <th>App Name</th>
      <th>App Version</th>
      <th>Environment</th>
      <th>Branch</th>
      <th>Node IP</th>
      <th>Updated Date</th>
      <th>Updated By</th>
    </tr>
    {{range .OneApp}}
    <tr>
        <td>{{ .ID }}</td>
        <td>{{ .AppName }}</td>
        <td>{{ .AppVersion }}</td>
        <td>{{ .Environment }}</td>
        <td>{{ .Branch }}</td>
        <td>{{ .IP }}</td>
        <td>{{ .UpdateDate }}</td>
        <td>{{ .UpdateBy }}</td>
    </tr>
    {{end}}
</table>
</body>
</html>`

// DisplayHTMLDev - Display HTML
func (a *App) DisplayHTMLDev(w http.ResponseWriter, r *http.Request) {

	var dev PageData

	dev, err := a.preperHTML("dev")

	if err != nil {
		log.Printf("cannot get dev aplication from db: %v", err)
	}

	tmpl, _ := template.New("dev").Parse(htm)
	tmpl.Execute(w, dev)
}

// DisplayHTMLStg - Display HTML for stage
func (a *App) DisplayHTMLStg(w http.ResponseWriter, r *http.Request) {

	var stg PageData

	stg, err := a.preperHTML("stg")

	if err != nil {
		log.Printf("cannot get stg aplication from db: %v", err)
	}

	tmpl, _ := template.New("stg").Parse(htm)
	tmpl.Execute(w, stg)
}

func (a *App) preperHTML(env string) (PageData, error) {

	var tmp AppStatusStruct

	appdata := []AppStatusStruct{}

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

	return data, nil
}
