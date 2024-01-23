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
		BackupIntervalSeconds int `ini:"backup_interval_in_seconds"`
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

	log.SetPrefix(time.Now().UTC().Format(HHMMSS24h) + ": ")
	log.SetFlags(log.Lmsgprefix)

	// Read .ini
	inidata, err := ini.Load("settings.ini")
	checkErr(err)

	var config Config

	// Map .ini
	err = inidata.MapTo(&config)
	checkErr(err)

	log.Println("Save Location: " + config.Paths.SaveLocation)
	log.Println("Backup Location: " + config.Paths.BackupLocation)

	// Read directory
	files, err := os.ReadDir(config.Paths.SaveLocation)
	checkErr(err)

	var newestFile fs.DirEntry
	var newestTime int64 = 0

	// Find last modified file
	for _, file := range files {
		fileInfo, err := file.Info()
		checkErr(err)
		log.Println(fileInfo.Name(), fileInfo.ModTime().Unix())

		currTime := fileInfo.ModTime().Unix()
		if currTime > newestTime {
			newestTime = currTime
			newestFile = file
		}
	}

	log.Println("Newest file name is: " + newestFile.Name())

	// Read source file
	in, err := os.Open(config.Paths.SaveLocation + "\\" + newestFile.Name())
	checkErr((err))
	defer in.Close()

	// Setup copy file
	timestamp := strconv.FormatInt(time.Now().UTC().UTC().UnixNano(), 10)
	out, err := os.Create(config.Paths.BackupLocation + "\\" + timestamp + "_" + newestFile.Name())
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
