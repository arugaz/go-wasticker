# go-wasticker

## Getting Started
Required ffmpeg, libwebp

### Install

```bash
> go get github.com/arugaz/go-wasticker
```

### Usage

```go
package main

import (
  "os"
  "github.com/arugaz/go-wasticker"
)

func main() {
  buf, _ := os.ReadFile("filename")

  stick, _ := wasticker.NewSticker(buf).ToByte()

  os.WriteFile("output", stick, 0644)

  url := "https://filename.ext"

  stick2, _ := wasticker.NewStickerUrl(url).ToByte()

  os.WriteFile("output", stick2, 0644)
}
```

## To-Do
 - Binding Libvips
 
---

- [RLOTTIE](https://github.com/Samsung/rlottie)
- [FFMPEG](https://github.com/FFmpeg/FFmpeg)
