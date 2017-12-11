package api

import (
	"net/http"

	"github.com/JamesClonk/home-info-dashboard/lib/web"
)

func Error(rw http.ResponseWriter, err error) {
	web.Render().JSON(rw, http.StatusInternalServerError, err.Error())
}
