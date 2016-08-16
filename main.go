package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"./proxy"
	"github.com/gorilla/mux"
	"github.com/yvasiyarov/gorelic"
	//"fmt"
)

var cfg proxy.Config

func startProxy() {
	if cfg.Threads > 0 {
		runtime.GOMAXPROCS(cfg.Threads)
		log.Printf("Running with %v threads", cfg.Threads)
	} else {
		n := runtime.NumCPU()
		runtime.GOMAXPROCS(n)
		log.Printf("Running with default %v threads", n)
	}

	r := mux.NewRouter()
	s := proxy.NewEndpoint(&cfg)

	r.Handle("/", s)
	err := http.ListenAndServe(cfg.Proxy.Listen, r)
	if err != nil {
		log.Fatal(err)
	}
}


func startNewrelic() {
	if cfg.NewrelicEnabled {
		nr := gorelic.NewAgent()
		nr.Verbose = cfg.NewrelicVerbose
		nr.NewrelicLicense = cfg.NewrelicKey
		nr.NewrelicName = cfg.NewrelicName
		nr.Run()
	}
}

func readConfig(cfg *proxy.Config) {
	configFileName := "config.json"
	//fmt.Println(os.Args)
	if len(os.Args) > 1 {
		configFileName = os.Args[1]
	}
	configFileName, _ = filepath.Abs(configFileName)
	log.Printf("Loading config: %v", configFileName)

	configFile, err := os.Open(configFileName)
	if err != nil {
		log.Fatal("File error: ", err.Error())
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&cfg); err != nil {
		log.Fatal("Config error: ", err.Error())
	}
}

func main() {
	readConfig(&cfg)
	startNewrelic()
	startProxy()
}
