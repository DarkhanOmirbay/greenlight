package main

import (
	"net/http"
)

// Declare a handler which writes a plain-text response with information about the
// application status, operating environment and version.
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln(w, "status: available")
	//fmt.Fprintf(w, "environment: %s\n", app.config.env)
	//fmt.Fprintf(w, "version: %s\n", version)

	//js := `{"status":"available", "enviroment": %q, "version": %q}`
	//js = fmt.Sprintf(js, app.config.env, version)
	//
	//w.Header().Set("Content-Type", "application/json")
	//
	//w.Write([]byte(js))

	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"enviroment": app.config.env,
			"version":    version,
		},
	}
	//js, err := json.Marshal(data)
	//if err != nil {
	//	app.logger.Print(err)
	//	http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	//	return
	//}
	//js = append(js, '\n')
	//w.Header().Set("Content-Type", "application/json")
	//w.Write([]byte(js))
	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		//app.logger.Print(err)
		//http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		app.serverErrorResponse(w, r, err)
	}
}
