package service

import (
	"bufio"
	"os"
)

type FileWorker struct {
	file   *os.File
	reader *bufio.Reader
}

func (fw *FileWorker) WriteData(data []byte) (int, error) {
	return fw.file.Write(data)
}

func (fw *FileWorker) ReadData() ([]byte, error) {
	var b []byte

	_, err := fw.reader.Read(b)

	if err != nil {
		return nil, err
	}

	return b, nil
}

func NewFileWorker(path string) (*FileWorker, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		return nil, err
	}

	return &FileWorker{
		file:   file,
		reader: bufio.NewReader(file),
	}, nil
}
