package wasticker

import "regexp"

var (
	disableCache   bool
	rgxUrl         = regexp.MustCompile(`^(?:https?://)?(?:[^/.\s]+\.)?[-a-zA-Z0-9@:%._+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_+.~#?&/=]*)`)
	defaultOptions = Options{
		Author: "djodi",
		Pack:   "bedess",
	}
)
