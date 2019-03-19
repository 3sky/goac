package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

//GetApp - gets app from API bu ID
func (c *Configuration) GetApp(i int) {
	var client http.Client
	var a interface{}

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
	fmt.Println(a)
}
