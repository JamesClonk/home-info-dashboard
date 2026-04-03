package api

import (
	"log"
	"net/http"

	"github.com/JamesClonk/home-info-dashboard/lib/web"
)

func Error(rw http.ResponseWriter, err error) {
	log.Printf("error: %v", err)
	web.Render().JSON(rw, http.StatusInternalServerError, err.Error())
}
