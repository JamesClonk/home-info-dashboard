package web

import (
	"net/http"

	"github.com/JamesClonk/home-info-dashboard/lib/util"
	"github.com/gorilla/mux"
)

func GetLocation(req *http.Request) (string, string) {
	// first, try to read values from gorilla mux
	vars := mux.Vars(req)
	canton := vars["canton"]
	city := vars["city"]

	// then, parse the form and try to read the values from POST data
	req.ParseForm()
	if len(canton) == 0 {
		canton = req.Form.Get("canton")
	}
	if len(city) == 0 {
		city = req.Form.Get("city")
	}

	return util.GetDefaultLocation(canton, city)
}
