package processing

/*
Anything image related:
- resize
- convert PNG → JPEG
- grayscale
- split double pages
- crop borders
*/

import (
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"os"

	"golang.org/x/image/draw"
)

func Grayscale(src image.Image) *image.Gray {
	bounds := src.Bounds()
	dst := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			g := color.GrayModel.Convert(src.At(x, y))
			dst.Set(x, y, g)
		}
	}

	return dst
}

// Rotate90CW rotates an image 90° clockwise.
func Rotate90CW(src image.Image) image.Image {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()

	dst := image.NewRGBA(image.Rect(0, 0, h, w))

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dst.Set(h-1-y, x, src.At(b.Min.X+x, b.Min.Y+y))
		}
	}

	return dst
}

func Resize(img image.Image, targetW, targetH int) image.Image {
	// Rotate landscape pages to portrait.
	if img.Bounds().Dx() > img.Bounds().Dy() {
		img = Rotate90CW(img)
	}

	srcW := img.Bounds().Dx()
	srcH := img.Bounds().Dy()

	// Preserve aspect ratio.
	scale := min(
		float64(targetW)/float64(srcW),
		float64(targetH)/float64(srcH),
	)

	newW := int(float64(srcW) * scale)
	newH := int(float64(srcH) * scale)

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))

	draw.CatmullRom.Scale(
		dst,
		dst.Bounds(),
		img,
		img.Bounds(),
		draw.Over,
		nil,
	)

	return dst
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func LoadSourceImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func EncodePNG(w io.Writer, img image.Image) error {
	return png.Encode(w, img)
}
