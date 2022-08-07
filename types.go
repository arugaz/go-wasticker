package wasticker

import (
	"image"
	"image/gif"
)

type (
	Options struct {
		Author     string
		Pack       string
		Categories []string
	}

	WASticker interface {
		ToByte() ([]byte, error)
		ToFile(filename string) error
	}

	newSticker struct {
		data     *[]byte
		url      *string
		metadata *Options
	}

	tgsgif struct {
		gif        gif.GIF
		images     []image.Image
		prev_frame *image.RGBA
	}

	imageWriter interface {
		init(w uint, h uint)
		addFrame(image *image.RGBA, fps uint) error
		result() []byte
	}
)
