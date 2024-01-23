package main

import (
	"log"
	"os"
	"time"

	"gopkg.in/ini.v1"
)

const (
	// 24h hh:mm:ss: 14:23:20
	HHMMSS24h = "15:04:05"
)

type Config struct {
	Paths struct {
		SaveLocation   string `ini:"save_location"`
		BackupLocation string `ini:"backup_location"`
	} `ini:"paths"`
	Timers struct {
		BackupIntervalSeconds int `ini:"backup_interval_in_seconds"`
	} `ini:"timers"`
}

func main() {
	// Setup logger
	logger := log.New(os.Stdout, time.Now().UTC().Format(HHMMSS24h)+": ", log.Lmsgprefix)

	// Read .ini
	inidata, err := ini.Load("settings.ini")
	if err != nil {
		logger.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	var config Config

	// Map .ini
	err = inidata.MapTo(&config)
	if err != nil {
		logger.Printf("Fail to map file %v", err)
		os.Exit(1)
	}

	logger.Println("Save Location: " + config.Paths.SaveLocation)
	logger.Println("Backup Location: " + config.Paths.BackupLocation)
}
