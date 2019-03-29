package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

/**
const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\n\033[1;31m%s %d\033[0m\n"
	DebugColor   = "\033[0;36m%s\033[0m"
)
**/
type AppStatusStruct struct {
	ID          int    `json:"id, omitempty"`
	AppName     string `json:"app_name"`
	AppVersion  string `json:"app_version"`
	Environment string `json:"environment"`
	Branch      string `json:"branch"`
	IP          string `json:"ip"`
	UpdateDate  string `json:"update_date"`
	UpdateBy    string `json:"update_by"`
}

var envs = map[string]string{
	"dev": "development",
	"stg": "stage",
}

var action = map[string]string{
	"get":     "gettin app by ID",
	"search":  "search app by name and env",
	"add":     "add new app",
	"update":  "update app",
	"promote": "promote app to next env",
	"delete":  "delete app",
}

//GetApp - gets app from API bu ID
func (c *Configuration) GetApp(i int) (AppStatusStruct, error) {
	var client http.Client
	var a AppStatusStruct

	url := fmt.Sprintf("http://%s:%d/api/app/%d", c.Server.IP, c.Server.Port, i)
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(c.Creditional.User, c.Creditional.Password)
	if err != nil {
		log.Printf("Error here", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error here", err)
	}

	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		fmt.Printf("Error while decode in GetApp: %v", err)
	}
	if len(a.AppName) != 0 {
		return a, nil
	} else {
		return AppStatusStruct{}, fmt.Errorf("There is no app with ID: %d", a.ID)
	}
}

//GetAppByName - Get app information with name and env
func (c *Configuration) GetAppByName(appPtr, environmentPtr string) (AppStatusStruct, error) {
	var client http.Client
	var a AppStatusStruct

	APP := &AppStatusStruct{
		AppName:     appPtr,
		Environment: environmentPtr,
	}

	payload, err := json.Marshal(APP)
	if err != nil {
		log.Printf("Error while marshal data", err)
	}

	url := fmt.Sprintf("http://%s:%d/api", c.Server.IP, c.Server.Port)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.SetBasicAuth(c.Creditional.User, c.Creditional.Password)
	if err != nil {
		log.Printf("Error here", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error here", err)
	}

	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		log.Printf("Error while decode in GetAppByName: %v", err)
	}

	if len(a.AppName) != 0 {
		return a, nil
	} else {
		return AppStatusStruct{}, fmt.Errorf("\nThere is no app with name %s on %s environment\n", appPtr, environmentPtr)
	}

}

//AddApp - add brand new app
func (c *Configuration) AddApp(appPtr, IPPtr, versionPtr, updaterPtr, environmentPtr, branchPtr string) (AppStatusStruct, error) {

	var client http.Client

	APP := &AppStatusStruct{
		AppName:     appPtr,
		AppVersion:  versionPtr,
		Environment: environmentPtr,
		Branch:      branchPtr,
		IP:          IPPtr,
		UpdateBy:    updaterPtr,
	}

	if _, ok := envs[environmentPtr]; !ok {
		fmt.Printf("\nWrong environment, You can use %s \n\n", getKeyFromMap(envs))
		os.Exit(0)
	}

	url := fmt.Sprintf("http://%s:%d/api/app/new", c.Server.IP, c.Server.Port)

	payload, err := json.Marshal(APP)
	if err != nil {
		fmt.Printf("Error while marshall in Insert App: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.SetBasicAuth(c.Creditional.User, c.Creditional.Password)
	if err != nil {
		log.Printf("Error here", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error here", err)
	}

	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	//TODO Insert sholud also make update

	if strings.Contains(buf.String(), "exits") {
		return AppStatusStruct{}, fmt.Errorf("\nThis app already exits on this environment !\n\n")
	} else {
		return c.GetAppByName(appPtr, environmentPtr)
	}

}

//DeleteApp - delete app by ID or new & env
func (c *Configuration) DeleteApp(appIDPtr int, appPtr, environmentPtr string) error {

	if !(appIDPtr == 0) {
		c.deleteAppByID(appIDPtr)
	} else if !(len(appPtr) == 0 || len(environmentPtr) == 0) {
		app, err := c.GetAppByName(appPtr, environmentPtr)
		if err != nil {
			return err
		}
		err = c.deleteAppByID(app.ID)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("ID or Name and Environment are empty !")
	}

	return nil
}

//TODO add promote
func (c *Configuration) PromoteApp(appIDPtr int, appPtr, environmentPtr string) (AppStatusStruct, error) {

	if !(appIDPtr == 0) {
		return c.promoteAppByID(appIDPtr)
	} else if !(len(appPtr) == 0 || len(environmentPtr) == 0) {
		base_app, err := c.GetAppByName(appPtr, environmentPtr)
		if err != nil {
			return AppStatusStruct{}, err
		}
		app, err := c.promoteAppByID(base_app.ID)
		if err != nil {
			return AppStatusStruct{}, err
		}

		return app, nil
	} else {
		return AppStatusStruct{}, fmt.Errorf("ID or Name and Environment are empty !")
	}
}

func (c *Configuration) promoteAppByID(i int) (AppStatusStruct, error) {

	var client http.Client

	var a AppStatusStruct
	url := fmt.Sprintf("http://%s:%d/api/app/%d", c.Server.IP, c.Server.Port, i)

	getReq, err := http.NewRequest("GET", url, nil)
	getReq.SetBasicAuth(c.Creditional.User, c.Creditional.Password)
	if err != nil {
		log.Printf("Error here", err)
	}

	getResp, err := client.Do(getReq)
	if err != nil {
		log.Printf("Error here", err)
	}

	defer getResp.Body.Close()
	err = json.NewDecoder(getResp.Body).Decode(&a)
	if err != nil {
		fmt.Printf("Error while decode in json : %v", err)
	}

	switch a.Environment {
	case "dev":
		a.Environment = "stg"
	case "stg":
		a.Environment = "preprod"
	case "preprod":
		a.Environment = "prod"
	}

	payload, err := json.Marshal(a)
	if err != nil {
		fmt.Printf("Error while marshall in promoteAppByID: %v", err)
	}

	url_new := fmt.Sprintf("http://%s:%d/api/app/new", c.Server.IP, c.Server.Port)
	postReq, err := http.NewRequest("POST", url_new, bytes.NewBuffer(payload))
	postReq.SetBasicAuth(c.Creditional.User, c.Creditional.Password)
	if err != nil {
		log.Printf("Error here", err)
	}

	postResp, err := client.Do(postReq)
	if err != nil {
		fmt.Printf("Error here", err)
	}
	defer postResp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(postResp.Body)
	fmt.Println(buf.String())
	if strings.Contains(buf.String(), "exits") {
		return AppStatusStruct{}, fmt.Errorf("\nThis app already exits on this environment !\n\n")
	} else {
		fmt.Println(a.AppName, a.Environment)
		return c.GetAppByName(a.AppName, a.Environment)
	}

}

func (c *Configuration) deleteAppByID(i int) error {

	var client http.Client

	url := fmt.Sprintf("http://%s:%d/api/app/%d", c.Server.IP, c.Server.Port, i)

	req, err := http.NewRequest("DELETE", url, nil)
	req.SetBasicAuth(c.Creditional.User, c.Creditional.Password)
	if err != nil {
		return fmt.Errorf("Error while creating DELETE request: %v \n", err)
	}

	_, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("Error while DO DELETE request: %v \n", err)
	}

	return nil

}
func getKeyFromMap(m map[string]string) []string {

	e := make([]string, 0, len(m))
	for k := range m {
		e = append(e, k)
	}
	return e
}

func (i *AppStatusStruct) Apssp() {
	i.Environment = "stg"
}
