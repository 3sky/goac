package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

//Insert - struct with data to insert
type Insert struct {
	ID          int
	IP          string
	AppName     string
	AppVersion  string
	Environment string
	Branch      string
	UpdateBy    string
}

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

var envs = map[string]string{
	"":    "empty",
	"dev": "development",
	"stg": "stage",
}

func main() {

	appIDPtr := flag.Int("id", 1, "application/component ID")
	appPtr := flag.String("app", "", "application/component name")
	IPPtr := flag.String("ip", "", "aplication IP")
	versionPtr := flag.String("ver", "", "version of application/component")
	updaterPtr := flag.String("updater", "", "person who insert this row")
	environmentPtr := flag.String("env", "", "application's environment")
	branchPtr := flag.String("branch", "", "application's branch")
	promotePtr := flag.Bool("p", false, "Promoting app to next environment")

	flag.Parse()

	cfg := LoadConfiguration(".creds")

	if _, ok := envs[*environmentPtr]; !ok {
		fmt.Println("#----------------#\nWrong environment, You can use", getKeyFromMap(envs))
		os.Exit(0)
	}

	ins := &Insert{
		ID:          *appIDPtr,
		IP:          *IPPtr,
		AppName:     *appPtr,
		AppVersion:  *versionPtr,
		Environment: *environmentPtr,
		Branch:      *branchPtr,
		UpdateBy:    *updaterPtr,
	}

	//app env from DB, not from cmd
	if *promotePtr && *environmentPtr == "dev" {
		ins.promoteApp()
	}

	cfg.InsertApp(ins)
	//cfg.GetApp(*appIDPtr)
}

func (i *Insert) promoteApp() {
	i.Environment = "stg"
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

func getKeyFromMap(m map[string]string) []string {

	e := make([]string, 0, len(m))
	for k := range m {
		e = append(e, k)
	}
	return e
}
