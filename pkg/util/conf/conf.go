package conf

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
)

func ChangeConfSite(site string, confPath string) (err error) {
	file, err := os.OpenFile(confPath, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	rd := bufio.NewReader(file)
	var line []byte
	var cursor int
	for err == nil {
		line, err = rd.ReadBytes('\n')
		if bytes.HasPrefix(line, []byte{'s', 'i', 't', 'e'}) {
			cursor += 6
			goto getCursor
		}
		cursor += len(line)
	}

	if err == io.EOF {
		return errors.New("not found site field")
	}

getCursor:
	_, err = file.WriteAt([]byte(site), int64(cursor))
	if err != nil {
		return
	}
	return file.Truncate(int64(cursor) + int64(len(site)))
}
