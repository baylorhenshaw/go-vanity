package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
)

type Config struct {
	Host    string   `json:"host"`
	Modules []Module `json:"modules"`
}

type Module struct {
	Name   string `json:"name"`
	Github string `json:"github"`
}

func main() {

	modules, config := registerModules()

	var port = os.Getenv("PORT")
	if port == "" {
		port = "3000"
		log.Info("Defaulting to port " + port)
	}

	log.Info("Server is running on port " + port)

	http.HandleFunc("/{module}", func(w http.ResponseWriter, r *http.Request) {
		if _, ok := modules[r.PathValue("module")]; ok {
			type ModulePageData struct {
				Module Module
				Host   string
			}
			tmp := template.Must(template.ParseFiles("templates/module.html"))
			tmp.Execute(w, ModulePageData{Module: modules[r.PathValue("module")], Host: config.Host})
		} else {
			w.WriteHeader(http.StatusNotFound)
			tmp := template.Must(template.ParseFiles("templates/404.html"))
			tmp.Execute(w, nil)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		log.Info("Received Request")

		type IndexPageData struct {
			Modules []Module
			Host    string
		}

		tmp := template.Must(template.ParseFiles("templates/index.html"))
		tmp.Execute(w, IndexPageData{Modules: config.Modules, Host: config.Host})
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// Register all modules from modules.json
func registerModules() (map[string]Module, Config) {

	configFile, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Print(err)
	}

	var config Config
	_ = json.Unmarshal([]byte(configFile), &config)

	var modules = make(map[string]Module)
	for i := 0; i < len(config.Modules); i++ {
		modules[config.Modules[i].Name] = config.Modules[i]
		log.Info("Registered module " + config.Modules[i].Name + " from " + config.Modules[i].Github)
	}

	return modules, config
}
