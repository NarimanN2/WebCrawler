package writer

import (
	"errors"
	"os"
)

type Writer interface {
	Write(path string, document []byte, append bool) error
	Mkdir(path string) error
}

type FileWriter struct {
}

func NewWriter() Writer {
	return FileWriter{
	}
}

func (fw FileWriter) Write(path string, document []byte, append bool) error {
	var file *os.File
	flag := os.O_CREATE | os.O_WRONLY

	if append {
		flag |= os.O_APPEND
	}

	file, err := os.OpenFile(path, flag, 0644)
	if err != nil {
		return errors.New("failed to write to writer")
	}

	defer file.Close()
	_, err = file.Write(document)
	if err != nil {
		return errors.New("failed to write to writer")
	}

	return nil
}

func (fw FileWriter) Mkdir(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}
