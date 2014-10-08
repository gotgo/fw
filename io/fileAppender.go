package io

import "os"

type FileAppender struct {
	FilePath string
	file     *os.File
}

func (fw *FileAppender) Open() error {
	if file, err := os.OpenFile(fw.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666); err != nil {
		return err
	} else {
		fw.file = file
	}

	return nil
}

func (fw *FileAppender) Write(b []byte) (n int, err error) {
	return fw.file.Write(b)
}

func (fw *FileAppender) Flush() error {
	return fw.file.Sync()
}

func (fw *FileAppender) Close() error {
	return fw.file.Close()
}
