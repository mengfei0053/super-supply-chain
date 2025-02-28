package configs

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
)

var DSN string
var DOMAIN string

var PORT string
var AuthKey string

var WEB_DAV_USER string
var WEB_DAV_PASSWORD string
var WEB_DAV_URL string
var ENVIRONMENT string

func LoadConfigFile() {
	var err error
	pwd, err := os.Getwd()
	fmt.Println("pwd", pwd)

	if err != nil {
		panic(err)
	}
	envDir := filepath.Join(filepath.Dir(pwd), "configs")

	envFile := filepath.Join(envDir, ".env")
	err = godotenv.Load(envFile)
	if err != nil {
		panic(err)
	}

	DSN = os.Getenv("DSN")
	DOMAIN = os.Getenv("DOMAIN")
	PORT = os.Getenv("PORT")
	AuthKey = os.Getenv("AUTH_KEY")
	WEB_DAV_USER = os.Getenv("WEB_DAV_USER")
	WEB_DAV_PASSWORD = os.Getenv("WEB_DAV_PASSWORD")
	WEB_DAV_URL = os.Getenv("WEB_DAV_URL")
	ENVIRONMENT = os.Getenv("ENVIRONMENT")
	fmt.Println("ENVIRONMENT", ENVIRONMENT)
}

func IsDev() bool {
	return ENVIRONMENT != "production"
}
