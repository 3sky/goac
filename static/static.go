package static

import (
	"p2go/db"
	//"p2go/api"
	"html/template"
	"net/http"
	// "fmt"
	// "time"
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
    // funcMap := template.FuncMap{
    //     "FormatDate": func(value time.Time) string {
    //         return fmt.Sprintf("%.s", value.Format("2006-01-02 15:04:05"))
    //     },
	// }
	
	var tmp db.StatusStruct

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


	tmpl := template.Must(template.ParseFiles("./static/hello.html"))
	tmpl.Execute(w, data)
}





// func main() {

// 	as1:= AS{APP_NAME: "Test 1", APP_VERSION: "1", UPDATE_BY: "1"}
// 	as2:= AS{APP_NAME: "Test 2", APP_VERSION: "2", UPDATE_BY: "2"}
 
// 	PagesData :=[]api.AppStatusStruct{}
	
// 	pages := PageData{PageTitle:"Test", Oneapp: PagesData} 

// 	pages.AddItem(as1)
// 	pages.AddItem(as2)

// 	fmt.Println(pages)
// }