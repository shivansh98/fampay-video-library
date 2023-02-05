package util

import "log"

func LogError(err error, section string) {
	log.Print("Got an error ", err, " in ", section)
}
