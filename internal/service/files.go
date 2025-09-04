package service

import (
	"bufio"
	"io"
	"os"
)

type FileWorker struct {
	file *os.File
	rw   *bufio.ReadWriter
}

func (fw *FileWorker) WriteData(data []byte) error {
	err := fw.file.Truncate(0)
	if err != nil {
		return err
	}

	_, err = fw.file.Seek(0, 0)
	if err != nil {
		return err
	}

	if _, err := fw.rw.Write(data); err != nil {
		return err
	}

	return fw.rw.Flush()
}

func (fw *FileWorker) ReadData() ([]byte, error) {
	_, err := fw.file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(fw.rw.Reader)

	return data, err
}

func (fw *FileWorker) Close() {
	fw.file.Close()
}

func NewFileWorker(path string, isSync bool) (*FileWorker, error) {
	var mask int
	if isSync {
		mask = os.O_RDWR | os.O_CREATE | os.O_SYNC
	} else {
		mask = os.O_RDWR | os.O_CREATE
	}
	file, err := os.OpenFile(path, mask, 0644)

	if err != nil {
		return nil, err
	}

	r := bufio.NewReader(file)
	w := bufio.NewWriter(file)

	return &FileWorker{
		file: file,
		rw:   bufio.NewReadWriter(r, w),
	}, nil
}
