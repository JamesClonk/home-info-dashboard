package util

import "github.com/JamesClonk/home-info-dashboard/lib/env"

func GetDefaultLocation(canton, city string) (string, string) {
	// try to read defaults from ENV, with reasonable defaults otherwise
	if len(canton) == 0 {
		canton = env.Get("DEFAULT_CANTON", "Bern")
	}
	if len(city) == 0 {
		city = env.Get("DEFAULT_CITY", "Bern")
	}
	return canton, city
}
