package wasticker

type (
	Options struct {
		Author     string
		Pack       string
		Categories []string
		Decrease   bool
	}

	WASticker interface {
		SetAuthor(author string)
		SetPack(pack string)
		SetCategories(categories []string)
		SetDecrease(decrease bool)
		ToByte() ([]byte, error)
		ToFile(filename string) error
	}

	newSticker struct {
		data     *[]byte
		metadata *Options
	}
)
