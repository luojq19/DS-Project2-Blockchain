package main

import "log"

func printError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
