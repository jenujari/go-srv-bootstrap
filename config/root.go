package config

import (
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	CONFIG_FOLDER = ".gosrv"
	CONFIG_DB     = "data.db"
)

var (
	dbPATH  string
	dirPATH string
	dbc     *gorm.DB
	Log      *log.Logger
)

func init() {
	// init log system
	Log = log.Default()
	Log.SetOutput(os.Stdout)

	// set default config folder in user default directory
	rootPath, err := os.UserHomeDir()
	if err != nil {
		Log.Panic("failed to get user directory")
	}

	dirPATH = filepath.Join(rootPath, CONFIG_FOLDER)
	dbPATH = filepath.Join(dirPATH, CONFIG_DB)
	os.MkdirAll(dirPATH, os.ModePerm)

	// setup db connection
	dbc, err = gorm.Open(sqlite.Open(dbPATH), &gorm.Config{})
	if err != nil {
		Log.Panic("failed to connect database")
	}
}

func GetDBC() *gorm.DB {
	return dbc
}

func GetUserDir() string {
	return dirPATH
}