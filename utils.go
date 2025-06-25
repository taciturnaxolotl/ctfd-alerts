package main

import (
	"io"
	"log"
	"os"
)

func setupLogging(debug bool) {
	if debug {
		debugLog = log.New(os.Stdout, "DEBUG: ", log.Ltime|log.Lmicroseconds)
	} else {
		debugLog = log.New(io.Discard, "", 0)
	}
}
