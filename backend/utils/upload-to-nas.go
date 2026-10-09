package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"super-supply-chain/configs"

	"github.com/studio-b12/gowebdav"
)

func isLocalUploadStore() bool {
	u := configs.WEB_DAV_URL
	return strings.HasPrefix(u, "file://") || strings.HasPrefix(u, "local://")
}

func localUploadStoreDir() (string, error) {
	u := configs.WEB_DAV_URL
	var dir string
	switch {
	case strings.HasPrefix(u, "file://"):
		dir = strings.TrimPrefix(u, "file://")
	case strings.HasPrefix(u, "local://"):
		dir = strings.TrimPrefix(u, "local://")
	default:
		return "", fmt.Errorf("not a local upload store URL: %s", u)
	}
	if dir == "" {
		return "", fmt.Errorf("local upload store path is empty")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}

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
	if isLocalUploadStore() {
		dir, err := localUploadStoreDir()
		if err != nil {
			return "", err
		}
		dest := filepath.Join(dir, fileName)
		if err := copyFile(filePath, dest); err != nil {
			return "", err
		}
		return configs.WEB_DAV_URL + "/" + fileName, nil
	}

	username := configs.WEB_DAV_USER
	password := configs.WEB_DAV_PASSWORD

	fileUrl := configs.WEB_DAV_URL + "/" + fileName

	err := uploadFileToWebDAV(filePath, configs.WEB_DAV_URL, username, password)
	if err != nil {
		return "", err
	}
	return fileUrl, nil
}
