package main

import (
	"net/http"
	"time"
)

var pathToTemplate = "./cmd/web/templates"

type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float64
	Data      map[string]any
	Flash     string
	Warning   string
	Error     string
	Now       time.Time
	// User *data.User
}

func (app *Config) render(w http.ResponseWriter, r *http.Request, t string, td *TemplateData) {

}
