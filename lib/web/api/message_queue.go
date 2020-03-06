package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/JamesClonk/home-info-dashboard/lib/database"
	"github.com/JamesClonk/home-info-dashboard/lib/web"
	"github.com/gorilla/mux"
)

func GetQueues(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		queues, err := hdb.GetQueues()
		if err != nil {
			Error(rw, err)
			return
		}
		web.Render().JSON(rw, http.StatusOK, queues)
	}
}

func GetAllMessages(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		messages, err := hdb.GetMessages()
		if err != nil {
			Error(rw, err)
			return
		}
		web.Render().JSON(rw, http.StatusOK, messages)
	}
}

func GetMessagesByQueue(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		queue := vars["queue"]

		messages, err := hdb.GetMessagesFromQueue(queue)
		if err != nil {
			Error(rw, err)
			return
		}
		web.Render().JSON(rw, http.StatusOK, messages)
	}
}

func GetMessage(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		var err error

		vars := mux.Vars(req)
		id := vars["id"]

		if len(id) > 0 {
			var messageId int64
			if len(id) > 0 {
				messageId, err = strconv.ParseInt(id, 10, 64)
				if err != nil {
					Error(rw, err)
					return
				}
			}

			message, err := hdb.GetMessageById(int(messageId))
			if err != nil {
				Error(rw, err)
				return
			}

			web.Render().JSON(rw, http.StatusOK, message)
			return
		}
		web.Render().JSON(rw, http.StatusNotFound, nil)
	}
}

func AddMessage(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			Error(rw, err)
			return
		}

		if len(req.Form.Get("queue")) == 0 {
			Error(rw, fmt.Errorf("queue missing!"))
			return
		}
		if len(req.Form.Get("message")) == 0 {
			Error(rw, fmt.Errorf("message missing!"))
			return
		}

		message := &database.Message{
			Queue:   req.Form.Get("queue"),
			Message: req.Form.Get("message"),
		}

		if err := hdb.InsertMessage(message); err != nil {
			Error(rw, err)
			return
		}
		web.Render().JSON(rw, http.StatusCreated, *message)
	}
}

func DeleteQueue(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		queue := vars["queue"]

		if len(queue) > 0 {
			if err := hdb.DeleteQueue(queue); err != nil {
				Error(rw, err)
				return
			}
			web.Render().JSON(rw, http.StatusNoContent, nil)
			return
		}
		web.Render().JSON(rw, http.StatusNotFound, nil)
	}
}

func DeleteMessage(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		var err error

		vars := mux.Vars(req)
		id := vars["id"]

		if len(id) > 0 {
			var messageId int64
			if len(id) > 0 {
				messageId, err = strconv.ParseInt(id, 10, 64)
				if err != nil {
					Error(rw, err)
					return
				}
			}

			if err := hdb.DeleteMessage(int(messageId)); err != nil {
				Error(rw, err)
				return
			}
			web.Render().JSON(rw, http.StatusNoContent, nil)
			return
		}
		web.Render().JSON(rw, http.StatusNotFound, nil)
	}
}
