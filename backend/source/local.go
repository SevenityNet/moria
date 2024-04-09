package source

import (
	"moria/config"
	"os"
)

type LocalSource struct {
	path string
}

func NewLocalSource() Source {
	return &LocalSource{
		path: config.GetSourceLocalPath(),
	}
}

func (s *LocalSource) Upload(file []byte, filename string) error {
	filePath := s.getFilePath(filename)

	err := os.WriteFile(filePath, file, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (s *LocalSource) Delete(filename string) error {
	filePath := s.getFilePath(filename)

	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}

func (s *LocalSource) Get(filename string) ([]byte, error) {
	filePath := s.getFilePath(filename)

	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (s *LocalSource) getFilePath(filename string) string {
	if _, err := os.Stat(s.path); os.IsNotExist(err) {
		os.MkdirAll(s.path, os.ModePerm)
	}

	return s.path + "/" + filename
}
