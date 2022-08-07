package wasticker

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// Convert video to webp
func (nS *newSticker) videoToWebp(ext string) ([]byte, error) {
	nameFile := fmt.Sprintf("%s%c%d", os.TempDir(), os.PathSeparator, time.Now().Unix())
	inFile := fmt.Sprintf("%s.%s", nameFile, ext)
	outFile := fmt.Sprintf("%s.webp", nameFile)

	defer os.Remove(inFile)
	defer os.Remove(outFile)
	err := os.WriteFile(inFile, *nS.data, 0644)
	if err != nil {
		return nil, err
	}

	err = exec.Command("ffmpeg", "-i", inFile, "-vf", "scale='min(512,iw)':min'(512,ih)':force_original_aspect_ratio=decrease,fps=30, pad=512:512:-1:-1:color=white@0.0, split [a][b]; [a] palettegen=reserve_transparent=on:transparency_color=ffffff [p]; [b][p] paletteuse", "-loop", "0", "-c:v", "libwebp", outFile, "-y").Run()
	if err != nil {
		return nil, err
	}

	return os.ReadFile(outFile)
}
