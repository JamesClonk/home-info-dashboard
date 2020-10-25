package api

import (
	"net/http"

	"github.com/JamesClonk/home-info-dashboard/lib/alerting/teams"
	"github.com/JamesClonk/home-info-dashboard/lib/web"
)

func TeamsStatus() func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		web.Render().JSON(rw, http.StatusOK, teams.Get())
	}
}

func TeamsMessage() func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		message := req.Form.Get("message")

		if len(message) == 0 {
			web.Render().JSON(rw, http.StatusBadRequest, map[string]string{"error": "message is empty!"})
			return
		}
		if err := teams.Get().Send(message); err != nil {
			Error(rw, err)
			return
		}

		web.Render().JSON(rw, http.StatusNoContent, nil)
	}
}
