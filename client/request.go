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
	ID          int    `json:"id"`
	AppName     string `json:"app_name"`
	AppVersion  string `json:"app_version"`
	Environment string `json:"environment"`
	Branch      string `json:"branch"`
	IP          string `json:"ip"`
	UpdateDate  string `json:"update_date"`
	UpdateBy    string `json:"update_by"`
}

var envs = map[string]string{
	"":    "empty",
	"dev": "development",
	"stg": "stage",
}

var action = map[string]string{
	"":        "empty",
	"get":     "gettin app by ID",
	"search":  "search app by name and env",
	"insert":  "insert app",
	"promote": "promote app to next env",
	"delete":  "delete app",
}

/**
	if _, ok := envs[*environmentPtr]; !ok {
		fmt.Printf("\nWrong environment, You can use %s \n\n", getKeyFromMap(envs))
		os.Exit(0)
	}



		if *promotePtr && *environmentPtr == "dev" {
			ins.promoteApp()
		}

**/

//GetApp - gets app from API bu ID
func (c *Configuration) GetApp(i int) error {
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
		a.prettyPrint()
	} else {
		fmt.Println("There is no app with ID:", a.ID)
	}
	return nil
}

//GetAppByName - Get app information with name and env
func (c *Configuration) GetAppByName(appPtr, environmentPtr string) error {
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
		fmt.Printf("Error while decode in TestDisplayAppByID: %v", err)
	}
	if len(a.AppName) != 0 {
		a.prettyPrint()
	} else {
		fmt.Printf("There is no app with name %s on %s environment\n", appPtr, environmentPtr)
	}

	return nil

}

func (c *Configuration) InsertApp(appPtr, IPPtr, versionPtr, updaterPtr, environmentPtr, branchPtr string) error {

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
		fmt.Printf("Error while marshall in TestAddNewApp: %v", err)
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

	if strings.Contains(buf.String(), "exits") {
		fmt.Printf("\nThis app already exits on this environment !\n\n")
	} else {
		c.GetAppByName(appPtr, environmentPtr)
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

func (i *AppStatusStruct) promoteApp() {
	i.Environment = "stg"
}

func (a *AppStatusStruct) prettyPrint() {
	longestString := len(fmt.Sprintf("UpdateDate: ", a.UpdateDate))
	for i := 0; i < longestString; i++ {
		fmt.Printf("#")
	}
	fmt.Println("#")
	fmt.Println("ID: ", a.ID)
	fmt.Println("AppName:", a.AppName)
	fmt.Println("AppVersion: ", a.AppVersion)
	fmt.Println("Environment: ", a.Environment)
	fmt.Println("Branch: ", a.Branch)
	fmt.Println("UpdateDate: ", a.UpdateDate)
	fmt.Println("UpdateBy: ", a.UpdateBy)
	for i := 0; i < longestString; i++ {
		fmt.Printf("#")
	}
	fmt.Println("#")

}
