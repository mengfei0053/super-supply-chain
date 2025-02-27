package utils

import (
	"github.com/studio-b12/gowebdav"
	"os"
	"super-supply-chain/configs"
)

func uploadFileToWebDAV(filePath, webdavURL, username, password string) error {
	client := gowebdav.NewClient(webdavURL, username, password)

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	err = client.WriteStream("/"+fileInfo.Name(), file, fileInfo.Mode())
	if err != nil {
		return err
	}

	return nil
}

func UploadToNas(filePath string, fileName string) (string, error) {
	username := configs.WEB_DAV_USER
	password := configs.WEB_DAV_PASSWORD

	fileUrl := configs.WEB_DAV_URL + "/" + fileName

	err := uploadFileToWebDAV(filePath, configs.WEB_DAV_URL, username, password)
	if err != nil {
		return "", err
	} else {
		return fileUrl, nil
	}
}
