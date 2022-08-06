package wasticker

import (
	"io"
	"net/http"
	"os"
	"unsafe"
)

// Create New Sticker
func NewSticker(data []byte, opts ...Options) WASticker {
	opt := &defaultOptions

	if len(opts) > 0 {
		opt = validateOpts(opts[0])
	}

	return &newSticker{data: &data, metadata: opt}
}

// Set Author
func (nS *newSticker) SetAuthor(author string) {
	nS.metadata.Author = author
}

// Set Pack Name
func (nS *newSticker) SetPack(pack string) {
	nS.metadata.Pack = pack
}

// Set Categories Emoji
func (nS *newSticker) SetCategories(catgeories []string) {
	nS.metadata.Categories = catgeories
}

// Set Decrease Size
func (nS *newSticker) SetDecrease(decrease bool) {
	nS.metadata.Decrease = decrease
}

// To Bytes
func (nS *newSticker) ToByte() ([]byte, error) {
	err := nS.parse()
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// To File
func (nS *newSticker) ToFile(filename string) error {
	data, err := nS.ToByte()
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (nS *newSticker) build() ([]byte, error) {

	return nil, nil
}

func (nS *newSticker) parse() error {
	isUrl := rgxUrl.Match(*nS.data)
	if isUrl {
		resp, err := http.Get(*(*string)(unsafe.Pointer(nS.data)))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		*nS.data = data
	}
	return nil
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
