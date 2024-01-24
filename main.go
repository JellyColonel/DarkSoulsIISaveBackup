package main

import (
	"io"
	"io/fs"
	"log"
	"os"
	"strconv"
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
		BackupInterval int `ini:"backup_interval_minutes"`
	} `ini:"timers"`
}

// Setup logger
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func main() {

	log.SetFlags(log.Ltime)

	// Read .ini
	inidata, err := ini.Load("settings.ini")
	checkErr(err)

	var config Config

	// Map .ini
	err = inidata.MapTo(&config)
	checkErr(err)

	var newestFile fs.DirEntry
	var newestTime int64 = 0

	for {

		// Read directory
		files, err := os.ReadDir(config.Paths.SaveLocation)
		checkErr(err)

		// Find last modified file
		for _, file := range files {
			fileInfo, err := file.Info()
			checkErr(err)

			currTime := fileInfo.ModTime().Unix()
			if currTime > newestTime {
				newestTime = currTime
				newestFile = file
			}
		}

		timestamp := strconv.FormatInt(time.Now().UTC().UTC().UnixNano(), 10)
		var sourcePath = config.Paths.SaveLocation + "\\" + newestFile.Name()
		var destinationPath = config.Paths.BackupLocation + "\\" + timestamp + "_" + newestFile.Name()
		createBackup(sourcePath, destinationPath)
		log.Println("Backup created")
		time.Sleep(time.Minute * time.Duration(config.Timers.BackupInterval))
	}
}

func createBackup(sourcePath string, destinationPath string) {

	// Read source file
	in, err := os.Open(sourcePath)
	checkErr((err))
	defer in.Close()

	// Setup copy file
	out, err := os.Create(destinationPath)
	checkErr(err)
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err := io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
}
