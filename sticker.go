package wasticker

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"unsafe"

	"github.com/arugaz/filetype"
)

// Create New Sticker
func NewSticker(data []byte, opts ...Options) WASticker {
	opt := &defaultOptions

	if len(opts) > 0 {
		opt = validateOpts(opts[0])
	}

	return &newSticker{data: &data, metadata: opt}
}

// Create New Sticker
func NewStickerUrl(url string, opts ...Options) WASticker {
	opt := &defaultOptions

	if len(opts) > 0 {
		opt = validateOpts(opts[0])
	}

	return &newSticker{url: &url, metadata: opt}
}

// To Bytes
func (nS *newSticker) ToByte() ([]byte, error) {
	err := nS.parse()
	if err != nil {
		return nil, err
	}

	ext, err := nS.getMime()
	if err != nil {
		return nil, err
	}
	data, err := nS.build(ext)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return data, nil
}

// To File
func (nS *newSticker) ToFile(filename string) error {
	data, err := nS.ToByte()
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// Parse data
func (nS *newSticker) parse() error {
	if len(*nS.data) == 0 && *nS.url != "" {
		resp, err := http.Get(*(*string)(unsafe.Pointer(nS.data)))
		if err != nil {
			return err
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		resp.Body.Close()
		*nS.data = data
	}
	if len(*nS.data) == 0 {
		return errors.New("there is no data")
	}
	return nil
}

func (nS *newSticker) getMime() (string, error) {
	types, err := filetype.Match(*nS.data)
	if err != nil {
		return "", err
	}
	return types.Extension, nil
}

func (nS *newSticker) build(ext string) ([]byte, error) {
	if ext == "tgs" {
		return nS.tgsToWebp()
	}
	return nS.videoToWebp(ext)
}

// Validation Options
func validateOpts(opt Options) *Options {
	if opt.Author == "" {
		opt.Author = defaultOptions.Author
	}
	if opt.Pack == "" {
		opt.Pack = defaultOptions.Pack
	}
	return &opt
}
