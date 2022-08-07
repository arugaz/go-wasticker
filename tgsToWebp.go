package wasticker

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"math"
	"os"
	"os/exec"
	"time"

	"github.com/arugaz/go-rlottie"
)

func (nS *newSticker) tgsToWebp() ([]byte, error) {
	if !disableCache {
		disableCache = true
		rlottie.LottieConfigureModelCacheSize(0)
	}
	uncompressed, err := tgsUnzip(*nS.data)
	if err != nil {
		return nil, err
	}
	animation := rlottie.LottieAnimationFromData(string(uncompressed[:]), "", "")
	if animation == nil {
		return nil, errors.New("failed to import lottie animation data")
	}

	w, h := rlottie.LottieAnimationGetSize(animation)
	w = uint(float32(w) * 1.)
	h = uint(float32(h) * 1.)

	frame_rate := rlottie.LottieAnimationGetFramerate(animation)
	frame_count := rlottie.LottieAnimationGetTotalframe(animation)
	duration := float32(frame_count) / float32(frame_rate)
	desired_framerate := float32(math.Min(30., frame_rate))

	if desired_framerate > 50. {
		desired_framerate = 50.
	}
	step := 1. / desired_framerate

	writer := newImageWriter(w, h)
	if writer == nil {
		return nil, errors.New("failed create imagewriter")
	}

	var i float32
	for i = 0.; i < duration; i += step {
		frame := rlottie.LottieAnimationGetFrameAtPos(animation, i/duration)
		buf := make([]byte, w*h*4)
		rlottie.LottieAnimationRender(animation, frame, buf, w, h, w*4)
		m := imageFromBuffer(buf, w, h)
		err := writer.addFrame(m, uint(desired_framerate))
		if err != nil {
			return nil, errors.New("Failed to add frame:" + err.Error())
		}
	}

	rlottie.LottieAnimationDestroy(animation)
	nameFile := fmt.Sprintf("%s%c%d", os.TempDir(), os.PathSeparator, time.Now().Unix())
	inFile := fmt.Sprintf("%s.gif", nameFile)
	outFile := fmt.Sprintf("%s.webp", nameFile)
	defer os.Remove(inFile)
	defer os.Remove(outFile)
	err = os.WriteFile(inFile, writer.result(), 0644)
	if err != nil {
		return nil, err
	}

	err = exec.Command("ffmpeg", "-i", inFile, "-loop", "0", "-c:v", "libwebp_anim", "-lossless", "1", outFile, "-y").Run()
	if err != nil {
		return nil, err
	}

	return os.ReadFile(outFile)
}

func (w_g *tgsgif) init(w uint, h uint) {
	w_g.gif.Config.Width = int(w)
	w_g.gif.Config.Height = int(h)
}

func (w_g *tgsgif) addFrame(image *image.RGBA, fps uint) error {
	var fps_int = int(1.0 / float32(fps) * 100.)
	if w_g.prev_frame != nil && sameImage(w_g.prev_frame, image) {
		w_g.gif.Delay[len(w_g.gif.Delay)-1] += fps_int
		return nil
	}
	w_g.gif.Image = append(w_g.gif.Image, nil)
	w_g.gif.Delay = append(w_g.gif.Delay, fps_int)
	w_g.gif.Disposal = append(w_g.gif.Disposal, gif.DisposalBackground)
	w_g.images = append(w_g.images, image)
	w_g.prev_frame = image
	return nil
}

func (w_g *tgsgif) result() []byte {
	q := medianCutQuantizer{mode, nil, false}
	p := q.quantizeMultiple(make([]color.Color, 0, 256), w_g.images)
	var trans_idx uint8 = 0
	if q.reserveTransparent {
		trans_idx = uint8(len(p))
	}
	var id_map = make(map[uint32]uint8)
	for i, img := range w_g.images {
		pi := image.NewPaletted(img.Bounds(), p)
		for y := 0; y < img.Bounds().Dy(); y++ {
			for x := 0; x < img.Bounds().Dx(); x++ {
				c := img.At(x, y)
				cr, cg, cb, ca := c.RGBA()
				cid := (cr>>8)<<16 | cg | (cb >> 8)
				if q.reserveTransparent && ca == 0 {
					pi.Pix[pi.PixOffset(x, y)] = trans_idx
				} else if val, ok := id_map[cid]; ok {
					pi.Pix[pi.PixOffset(x, y)] = val
				} else {
					val := uint8(p.Index(c))
					pi.Pix[pi.PixOffset(x, y)] = val
					id_map[cid] = val
				}
			}
		}
		w_g.gif.Image[i] = pi
	}
	if q.reserveTransparent {
		p = append(p, color.RGBA{0, 0, 0, 0})
	}
	for _, img := range w_g.gif.Image {
		img.Palette = p
	}
	w_g.gif.Config.ColorModel = p
	var data []byte
	w := bytes.NewBuffer(data)
	err := gif.EncodeAll(w, &w_g.gif)
	if err != nil {
		return nil
	}
	return w.Bytes()
}

func sameImage(a *image.RGBA, b *image.RGBA) bool {
	if len(a.Pix) != len(b.Pix) {
		return false
	}
	for i, v := range a.Pix {
		if v != b.Pix[i] {
			return false
		}
	}
	return true
}

func imageFromBuffer(p []byte, w uint, h uint) *image.RGBA {
	for i := 0; i < len(p); i += 4 {
		p[i+0], p[i+2] = p[i+2], p[i+0]
	}
	m := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	m.Pix = p
	m.Stride = int(w) * 4
	return m
}

func newImageWriter(w uint, h uint) imageWriter {
	writer := &tgsgif{}
	writer.init(w, h)
	return writer
}
