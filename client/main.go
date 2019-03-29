package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

//Configuration - struct with conf
type Configuration struct {
	Creditional `json:"creditional"`
	Server      `json:"server"`
}

//Creditional - struct with user's creditionals
type Creditional struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

//Server - struct with server info
type Server struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

var app AppStatusStruct

func main() {

	actionPtr := flag.String("action", "get", "what You want to do?")
	appIDPtr := flag.Int("id", 0, "application/component ID")
	appPtr := flag.String("app", "", "application/component name")
	IPPtr := flag.String("ip", "", "aplication IP")
	versionPtr := flag.String("ver", "", "version of application/component")
	updaterPtr := flag.String("updater", "", "person who insert this row")
	environmentPtr := flag.String("env", "", "application's environment")
	branchPtr := flag.String("branch", "", "application's branch")

	flag.Parse()

	cfg := LoadConfiguration(".creds")

	switch *actionPtr {
	case "get":
		fmt.Println("I will get it!")
		app, err := cfg.GetApp(*appIDPtr)
		if err != nil {
			fmt.Println(err)
		} else {
			app.prettyPrint()
		}
	case "search":
		fmt.Println("I will search it!")
		if len(*appPtr) != 0 && len(*environmentPtr) != 0 {
			app, err := cfg.GetAppByName(*appPtr, *environmentPtr)
			if err != nil {
				fmt.Println(err)
			} else {
				app.prettyPrint()
			}
		} else {
			fmt.Println(" \nApp and environment are mandatory !")
		}

	case "add":
		fmt.Println("I will add it!")
		app, err := cfg.AddApp(*appPtr, *IPPtr, *versionPtr, *updaterPtr, *environmentPtr, *branchPtr)
		if err != nil {
			fmt.Println("Error while inserting new app", err)
		} else {
			app.prettyPrint()
		}
	case "update":
		fmt.Println("I will update it!")
	case "promote":
		fmt.Println("I will promote it!")
		app, err := cfg.PromoteApp(*appIDPtr, *appPtr, *environmentPtr)
		if err != nil {
			fmt.Println("Error while promote app", err)
		} else {
			app.prettyPrint()
		}
	case "delete":
		fmt.Println("I will delete it!")
		err := cfg.DeleteApp(*appIDPtr, *appPtr, *environmentPtr)
		if err != nil {
			fmt.Println("Error while deleting app", err)
		} else {
			fmt.Println("App deleted ")
		}
	default:
		fmt.Println("I will do nothing! Valid action is:", getKeyFromMap(action))
	}
}

//LoadConfiguration - just load configuration file
func LoadConfiguration(file string) Configuration {
	var config Configuration
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

//PrettyPrint - just print app struct
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
