package web

import (
	"net/http"

	"github.com/JamesClonk/home-info-dashboard/lib/util"
	"github.com/gorilla/mux"
)

func GetLocation(req *http.Request) (string, string, string) {
	// first, try to read values from gorilla mux
	vars := mux.Vars(req)
	lat := vars["latitude"]
	lon := vars["longitude"]
	alt := vars["altitude"]

	// then, parse the form and try to read the values from POST data
	_ = req.ParseForm()
	if len(lat) == 0 {
		lat = req.Form.Get("latitude")
	}
	if len(lon) == 0 {
		lon = req.Form.Get("longitude")
	}
	if len(alt) == 0 {
		alt = req.Form.Get("altitude")
	}

	return util.GetDefaultLocation(lat, lon, alt)
}
