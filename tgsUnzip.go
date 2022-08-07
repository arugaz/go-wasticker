package wasticker

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
)

func tgsUnzip(data []byte) ([]byte, error) {
	z, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, errors.New("Failed to create gzip reader:" + err.Error())
	}
	unz, err := io.ReadAll(z)
	if err != nil {
		return nil, errors.New("Failed to read gzip archive:" + err.Error())
	}
	z.Close()
	return unz, nil
}
