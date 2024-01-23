package main

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
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
	log.SetFlags(log.Ltime)

	// Read .ini
	inidata, err := ini.Load("settings.ini")
	if err != nil {
		log.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	var config Config

	// Map .ini
	err = inidata.MapTo(&config)
	if err != nil {
		log.Printf("Fail to map file %v", err)
		os.Exit(1)
	}

	log.Println("Save Location: " + config.Paths.SaveLocation)
	log.Println("Backup Location: " + config.Paths.BackupLocation)
}
