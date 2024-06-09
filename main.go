package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

type Config struct {
	Host      string
	Modules   []Module `json:"modules"`
	Analytics Analytics
}

type Module struct {
	Name   string `json:"name"`
	Github string `json:"github"`
}

type Analytics struct {
	Enabled bool   `json:"enabled"`
	Url     string `json:"url"`
	Id      string `json:"id"`
}

func main() {

	loadEnvironment()

	modules, config := registerModules()

	var port = os.Getenv("PORT")
	if port == "" {
		port = "3000"
		log.Warn("Could not find PORT in environment variable, defaulting to 3000")
	}

	log.Info("Server is started on http://" + config.Host + ":" + port + "/")

	http.HandleFunc("/{module}", func(w http.ResponseWriter, r *http.Request) {
		if _, ok := modules[r.PathValue("module")]; ok {
			type ModulePageData struct {
				Module    Module
				Host      string
				Analytics Analytics
			}
			tmp := template.Must(template.ParseFiles("templates/module.tmpl"))
			tmp.Execute(w, ModulePageData{Module: modules[r.PathValue("module")], Host: config.Host, Analytics: config.Analytics})
			log.Info("Redirecting user to " + modules[r.PathValue("module")].Github)
		} else {
			w.WriteHeader(http.StatusNotFound)
			tmp := template.Must(template.ParseFiles("templates/404.tmpl"))
			tmp.Execute(w, config.Analytics)
			log.Warn("Module not found for /" + r.PathValue("module"))
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		log.Info("Received request for /")

		type IndexPageData struct {
			Modules   []Module
			Host      string
			Analytics Analytics
		}

		tmp := template.Must(template.ParseFiles("templates/index.tmpl"))
		tmp.Execute(w, IndexPageData{Modules: config.Modules, Host: config.Host, Analytics: config.Analytics})
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// Register all modules from modules.json
func registerModules() (map[string]Module, Config) {

	// Load Modules File
	configFile, err := os.ReadFile(os.Getenv("MODULES_FILE"))
	if err != nil {
		log.Error("Could not load modules file!")
	}

	// Parse Modules File
	var config Config
	_ = json.Unmarshal([]byte(configFile), &config)

	var modules = make(map[string]Module)
	for i := 0; i < len(config.Modules); i++ {
		modules[config.Modules[i].Name] = config.Modules[i]
		log.Info("Registered module " + config.Modules[i].Name + " from " + config.Modules[i].Github)
	}

	// Set Host
	if os.Getenv("HOST") != "" {
		config.Host = os.Getenv("HOST")
	} else {
		config.Host = "localhost"
		log.Warn("Could not find HOST in environment variable, defaulting to localhost")
	}

	// Enable Umami Analytics
	if os.Getenv("UMAMI_URL") != "" {
		if os.Getenv("UMAMI_URL") != "" {
			config.Analytics.Enabled = true
			config.Analytics.Url = os.Getenv("UMAMI_URL")
			config.Analytics.Id = os.Getenv("UMAMI_ID")
			log.Info("Enabled Umami Analytics")
		}
	}

	return modules, config
}

func loadEnvironment() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
