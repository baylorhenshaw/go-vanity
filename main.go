package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
)

type Module struct {
	Name   string `json:"name"`
	Author string `json:"author"`
	Github string `json:"github"`
}

func main() {

	modules := registerModules()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	log.Info("Server is running on port " + port)

	http.HandleFunc("/{module}", func(w http.ResponseWriter, r *http.Request) {
		if _, ok := modules[r.PathValue("module")]; ok {
			tmp := template.Must(template.ParseFiles("templates/module.html"))
			tmp.Execute(w, modules[r.PathValue("module")])
		} else {
			http.NotFound(w, r)
		}
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// Register all modules from modules.json
func registerModules() map[string]Module {

	moduleFile, err := os.ReadFile("modules.json")
	if err != nil {
		fmt.Print(err)
	}

	var moduleSlice []Module
	_ = json.Unmarshal([]byte(moduleFile), &moduleSlice)

	var modules = make(map[string]Module)
	for i := 0; i < len(moduleSlice); i++ {
		modules[moduleSlice[i].Name] = moduleSlice[i]
		log.Info("Registered module: " + moduleSlice[i].Name + " by " + moduleSlice[i].Author + " from " + moduleSlice[i].Github)
	}

	return modules
}
