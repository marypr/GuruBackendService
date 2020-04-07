package main

import (
	"github.com/marypr/GuruBackendService/src"
	"log"
	"os"
)

func main() {
	file := setupLogFile()
	defer file.Close()
	src.Connect()
	src.Start()
}

func setupLogFile() *os.File {
	logFile, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
	println("Error logs will be in the log.txt.")
	log.Println("Recording of the log file has started...")
	return logFile
}
