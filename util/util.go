package util

import "log"

func LogError(err error, section string) {
	log.Printf("Got an error %s in %s", err.Error(), section)
}
