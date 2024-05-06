package logger

import (
	configs "helpdesk_backend/config"
	"log"
	"os"
	"strconv"
	"time"
)

// Logger central logger
var Logger *log.Logger

// Logger out file
var outfile *os.File

// LogFileChanger : Logger setup and file changer
func LogFileChanger() {

	outfile, _ = os.Create(configs.LOG_PATH)

	Logger = log.New(outfile, "", 0)
	currTime := time.Now()
	time.Sleep(time.Duration(24-currTime.Hour()) * time.Hour)
	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		y, m, d := time.Now().Date()
		outfile, err := os.Create("logger-" + strconv.Itoa(d) + "-" + strconv.Itoa(int(m)) + "-" + strconv.Itoa(y) + ".txt")
		if err == nil {
			Logger = log.New(outfile, "", 0)
		} else {
			Logger.Println("error creating new file", err)
		}
	}
}
