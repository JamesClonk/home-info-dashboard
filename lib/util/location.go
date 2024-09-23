package util

import "github.com/JamesClonk/home-info-dashboard/lib/env"

func GetDefaultLocation(lat, lon, alt string) (string, string, string) {
	// try to read defaults from ENV, with reasonable defaults otherwise
	if len(lat) == 0 {
		lat = env.Get("DEFAULT_LATITUDE", "47.02115")
	}
	if len(lon) == 0 {
		lon = env.Get("DEFAULT_LONGITUDE", "7.44914")
	}
	if len(alt) == 0 {
		alt = env.Get("DEFAULT_ALTITUDE", "555")
	}
	return lat, lon, alt
}
