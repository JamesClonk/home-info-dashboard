package api

import (
	"net/http"

	"github.com/JamesClonk/home-info-dashboard/lib/slack"
	"github.com/JamesClonk/home-info-dashboard/lib/web"
)

func SlackStatus() func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		web.Render().JSON(rw, http.StatusOK, slack.Get())
	}
}

func SlackMessage() func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		message := req.Form.Get("message")

		if len(message) == 0 {
			web.Render().JSON(rw, http.StatusBadRequest, map[string]string{"error": "message is empty!"})
			return
		}
		if err := slack.Get().Send(message); err != nil {
			Error(rw, err)
			return
		}

		web.Render().JSON(rw, http.StatusNoContent, nil)
	}
}
