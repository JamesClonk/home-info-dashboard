package util

import "github.com/JamesClonk/home-info-dashboard/lib/env"

func GetUserAndPassword() (string, string) {
	return env.Get("WEATHERAPI_USERNAME", "jamesclonk"), env.MustGet("WEATHERAPI_PASSWORD")
}
