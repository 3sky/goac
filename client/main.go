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

//TODO add infoPrint function

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
		if *appIDPtr != 0 {
			app, err := cfg.GetApp(*appIDPtr)
			if err != nil {
				infoPrint(fmt.Sprintf("%v", err))
			} else {
				app.prettyPrint()
			}
		} else {
			infoPrint("App ID is mandatory !")
		}
	case "search":
		if len(*appPtr) != 0 && len(*environmentPtr) != 0 {
			app, err := cfg.GetAppByName(*appPtr, *environmentPtr)
			if err != nil {
				infoPrint(fmt.Sprintf("%v", err))
			} else {
				app.prettyPrint()
			}
		} else {
			infoPrint("App name and environment are mandatory !")
		}

	case "add":
		if len(*appPtr) != 0 && len(*environmentPtr) != 0 {
			app, err := cfg.AddApp(*appPtr, *IPPtr, *versionPtr, *updaterPtr, *environmentPtr, *branchPtr)
			if err != nil {
				infoPrint(fmt.Sprintf("Error while inserting new app %v", err))
			} else {
				app.prettyPrint()
			}
		} else {
			infoPrint("App name and environment are mandatory !")
		}
	case "update":
		if *appIDPtr != 0 {
			app, err := cfg.UpdateApp(*appIDPtr, *appPtr, *IPPtr, *versionPtr, *updaterPtr, *environmentPtr, *branchPtr)
			if err != nil {
				infoPrint(fmt.Sprintf("Error while update app %v", err))
			} else {
				app.prettyPrint()
			}
		} else {
			infoPrint("App ID is mandatory in update command !")
		}
	case "promote":
		//TODO bad loggin error info
		app, err := cfg.PromoteApp(*appIDPtr, *appPtr, *environmentPtr)
		if err != nil {
			infoPrint(fmt.Sprintf("Error while promote app %v", err))
		} else {
			app.prettyPrint()
		}
	case "delete":
		err := cfg.DeleteApp(*appIDPtr, *appPtr, *environmentPtr)
		if err != nil {
			infoPrint(fmt.Sprintf("Error while deleting app %v", err))
		} else {
			infoPrint("App deleted !")
		}
	default:
		infoPrint(fmt.Sprintf("I will do nothing! Valid action is: %s", getKeyFromMap(action)))
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
	longestString := len(fmt.Sprintf("UpdateDate: %s", a.UpdateDate))
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

func infoPrint(info string) {
	l := len(info) + 7
	for i := 0; i < l; i++ {
		fmt.Printf("=")
	}
	fmt.Printf("\n>>>  %s\n", info)
	for i := 0; i < l; i++ {
		fmt.Printf("=")
	}
	fmt.Println("")

}
