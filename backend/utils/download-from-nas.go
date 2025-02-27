package utils

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/studio-b12/gowebdav"
	"io"
	"log"
	"os"
	"path/filepath"
	"super-supply-chain/configs"
)

func DownloadFromNas(fileName string) (string, error) {
	uploadDir := GetUploadTmpDir()
	uuidFileName := uuid.New().String()
	extension := filepath.Ext(fileName)
	newFileName := uuidFileName + extension
	localFilePath := filepath.Join(uploadDir, newFileName)

	c := gowebdav.NewClient(configs.WEB_DAV_URL, configs.WEB_DAV_USER, configs.WEB_DAV_PASSWORD)
	reader, err := c.ReadStream(fileName)
	if err != nil {
		fmt.Println("Error reading file from NAS")
		log.Fatal(err)
		return "", err
	}

	file, err := os.Create(localFilePath)
	if err != nil {
		fmt.Println("Error creating file")
		log.Fatal(err)
		return "", err
	}
	defer file.Close()
	t, err := io.Copy(file, reader)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	fmt.Println(t)
	return localFilePath, nil

}
