package util

import "log"

func HandleError(err error, section string) {
	if err != nil {
		log.Print("Got an error ", err, " in ", section)
		panic(err)
	}
}

func LogError(err error, section string) {
	log.Print("Got an error ", err, " in ", section)
}
