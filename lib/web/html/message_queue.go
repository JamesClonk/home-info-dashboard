package html

import (
	"net/http"

	"github.com/JamesClonk/home-info-dashboard/lib/database"
	"github.com/JamesClonk/home-info-dashboard/lib/web"
)

func Messages(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		page := &Page{
			Title:  "Home Automation - Message Queue",
			Active: "message_queue",
		}

		// parse the form and try to read the values from POST data
		_ = req.ParseForm()
		queue := req.Form.Get("queue")
		message := req.Form.Get("message")
		if len(queue) > 0 && len(message) > 0 {
			if err := hdb.InsertMessage(&database.Message{
				Queue:   queue,
				Message: message,
			}); err != nil {
				Error(rw, err)
				return
			}
		}

		queues, err := hdb.GetQueues()
		if err != nil {
			Error(rw, err)
			return
		}

		messages, err := hdb.GetMessages()
		if err != nil {
			Error(rw, err)
			return
		}

		page.Content = struct {
			Queues   []*database.Queue
			Messages []*database.Message
		}{
			queues,
			messages,
		}
		_ = web.Render().HTML(rw, http.StatusOK, "messages", page)
	}
}
