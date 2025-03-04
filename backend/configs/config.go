package configs

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
)

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

	fmt.Println("ENVIRONMENT", os.Getenv("ENVIRONMENT"))
	fmt.Println("IsDev", IsDev())
	if IsDev() {

		envDir := filepath.Join(filepath.Dir(pwd), "configs")
		envFile := filepath.Join(envDir, ".env")
		err = godotenv.Load(envFile)
		if err != nil {
			panic(err)
		}
	}

	PORT = os.Getenv("PORT")
	AuthKey = "Authorization"
	WEB_DAV_USER = os.Getenv("UPLOAD_USER")
	WEB_DAV_PASSWORD = os.Getenv("UPLOAD_PASSWORD")
	WEB_DAV_URL = os.Getenv("UPLOAD_SERVER")
	ENVIRONMENT = os.Getenv("ENVIRONMENT")
	fmt.Println("ENVIRONMENT", ENVIRONMENT)
	fmt.Println("WEB_DAV_URL", WEB_DAV_URL)
	fmt.Println("WEB_DAV_USER", WEB_DAV_USER)
	fmt.Println("WEB_DAV_PASSWORD", WEB_DAV_PASSWORD)
	fmt.Println("PORT", PORT)
	
}

func IsDev() bool {
	return os.Getenv("ENVIRONMENT") != "production"
}
