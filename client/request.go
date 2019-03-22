package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type AppStatusStruct struct {
	ID          uint
	AppName     string `json:"app_name"`
	AppVersion  string `json:"app_version"`
	Environment string `json:"environment"`
	Branch      string `json:"branch"`
	UpdateDate  string `json:"update_date"`
	UpdateBy    string `json:"update_by"`
}

//GetApp - gets app from API bu ID
func (c *Configuration) GetApp(i int) {
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
		fmt.Printf("Error while decode in TestDisplayAppByID: %v", err)
	}
	if len(a.AppName) != 0 {
		a.prettyPrint()
	} else {
		fmt.Println("There is no app with ID:", a.ID)
	}
}

func (c *Configuration) InsertApp(i *Insert) {

	//var client http.Client

	url := fmt.Sprintf("http://%s:%d/api/app/new", c.Server.IP, c.Server.Port)

	jsonP1, err := json.Marshal(i)
	if err != nil {
		fmt.Printf("Error while marshall in TestAddNewApp: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonP1))
	req.SetBasicAuth(c.Creditional.User, c.Creditional.Password)
	if err != nil {
		log.Printf("Error here", err)
	}
	fmt.Println(i)
	fmt.Println(req)

	//resp, err := client.Do(req)
	//if err != nil {
	//	log.Printf("Error here", err)
	//}

}

func (a *AppStatusStruct) prettyPrint() {
	fmt.Println("ID: ", a.ID)
	fmt.Println("AppName:", a.AppName)
	fmt.Println("AppVersion: ", a.AppVersion)
	fmt.Println("Environment: ", a.Environment)
	fmt.Println("Branch: ", a.Branch)
	fmt.Println("UpdateDate: ", a.UpdateDate)
	fmt.Println("UpdateBy: ", a.UpdateBy)
}
