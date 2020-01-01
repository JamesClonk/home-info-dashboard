package util

import "github.com/JamesClonk/home-info-dashboard/lib/env"

func GetUserAndPassword() (string, string) {
	return env.Get("AUTH_USERNAME", "jamesclonk"), env.MustGet("AUTH_PASSWORD")
}
