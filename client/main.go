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

func main() {

	actionPtr := flag.String("action", "get", "what You want to do?")
	appIDPtr := flag.Int("id", 1, "application/component ID")
	appPtr := flag.String("app", "", "application/component name")
	IPPtr := flag.String("ip", "", "aplication IP")
	versionPtr := flag.String("ver", "", "version of application/component")
	updaterPtr := flag.String("updater", "", "person who insert this row")
	environmentPtr := flag.String("env", "", "application's environment")
	branchPtr := flag.String("branch", "", "application's branch")
	//promotePtr := flag.Bool("p", false, "Promoting app to next environment")

	flag.Parse()

	cfg := LoadConfiguration(".creds")

	switch *actionPtr {
	case "get":
		fmt.Println("I wiil get it!")
		err := cfg.GetApp(*appIDPtr)
		if err != nil {
			fmt.Println("Error while getting app", err)
		}
	case "search":
		fmt.Println("I will search it!")
		err := cfg.GetAppByName(*appPtr, *environmentPtr)
		if err != nil {
			fmt.Println("Error while searching app", err)
		}
	case "insert":
		fmt.Println("I will insert it!")
		err := cfg.InsertApp(*appPtr, *IPPtr, *versionPtr, *updaterPtr, *environmentPtr, *branchPtr)
		if err != nil {
			fmt.Println("Error while inserting new app", err)
		}
	case "promote":
		fmt.Println("I will promote it!")
		//cfg.PromoteApp()
	case "delete":
		fmt.Println("I will delete it!")
		//cfg.DeleteApp(*IPPtr)
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
