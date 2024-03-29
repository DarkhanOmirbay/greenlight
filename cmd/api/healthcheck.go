package main

import (
	"fmt"
	"net/http"
)

// Declare a handler which writes a plain-text response with information about the
// application status, operating environment and version.
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln(w, "status: available")
	//fmt.Fprintf(w, "environment: %s\n", app.config.env)
	//fmt.Fprintf(w, "version: %s\n", version)
	js := `{"status":"available", "enviroment": %q, "version": %q}`
	js = fmt.Sprintf(js, app.config.env, version)

	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(js))
}
