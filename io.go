package main

import (
	"errors"
	"log"
	"os"
	"path"

	"github.com/google/uuid"
)

const UPLOAD_FOLDER = ".uploads"
const TMP_FOLDER = ".tmp"

var (
	ErrFileNotFound = errors.New("file not found")
)

func initIO() {
	if _, err := os.Stat(UPLOAD_FOLDER); os.IsNotExist(err) {
		if err := os.Mkdir(UPLOAD_FOLDER, 0755); err != nil {
			panic(err)
		}
	}

	if _, err := os.Stat(TMP_FOLDER); os.IsNotExist(err) {
		if err := os.Mkdir(TMP_FOLDER, 0755); err != nil {
			panic(err)
		}
	}
}

func getFilePath(subfolder, filename string) string {
	return path.Join(UPLOAD_FOLDER, subfolder, filename)
}

func writeFileToDisk(subfolder, filename string, data []byte) error {
	fileDir := path.Join(UPLOAD_FOLDER, subfolder)
	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		if err := os.Mkdir(fileDir, 0755); err != nil {
			return err
		}
	}

	filePath := path.Join(fileDir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func deleteFileFromDisk(subfolder, filename string) error {
	filePath := getFilePath(subfolder, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return ErrFileNotFound
	}

	if err := os.Remove(filePath); err != nil {
		return err
	}

	if err := removeCachedFile(filename); err != nil {
		log.Println("Error removing cached file:", err)
	}

	return nil
}

func getTmpFilePath(ext string) string {
	return path.Join(TMP_FOLDER, uuid.NewString()+"."+ext)
}

func writeTmpFile(data []byte, ext string) (string, error) {
	path := getTmpFilePath(ext)

	file, err := os.Create(path)
	if err != nil {
		return "", err
	}

	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return "", err
	}

	return path, nil
}

func moveTmpToUploads(tmpPath, subfolder, filename string) error {
	fileDir := path.Join(UPLOAD_FOLDER, subfolder)
	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		if err := os.Mkdir(fileDir, 0755); err != nil {
			return err
		}
	}

	filePath := path.Join(fileDir, filename)
	if err := os.Rename(tmpPath, filePath); err != nil {
		return err
	}

	return nil
}

func deleteTmpFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ErrFileNotFound
	}

	if err := os.Remove(path); err != nil {
		return err
	}

	return nil
}
