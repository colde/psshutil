package fileHandling

import (
	"log"
	"os"
)

func ReadFromFile(f *os.File, size int64) ([]byte, error) {
	buf := make([]byte, size)
	_, err := f.Read(buf)

	if err != nil {
		log.Fatalln(err.Error())
		return nil, err
	}

	return buf, nil
}

func ReadHeader(f *os.File) ([]byte, []byte, error) {
	buf, err := ReadFromFile(f, 8)
	if err != nil {
		log.Fatalln(err.Error())
		return nil, nil, err
	}

	return buf[0:4], buf[4:8], nil
}
