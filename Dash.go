package dash

import (
	"encoding/hex"
	"github.com/nfnt/resize"
	"image"
	"image/color"
)

// converts image to grayscale
func grayscale(img image.Image) image.Image {
	gray := image.NewGray(img.Bounds())
	for x := 0 ; x < img.Bounds().Size().X; x++ {
		for y := 0; y < img.Bounds().Size().Y; y++ {
			col := color.GrayModel.Convert(img.At(x, y))
			gray.Set(x, y, col)
		}
	}

	return gray
}

// generates a dhash for the given image and hash length.
// hash is based on the image data and not it's container.
func Hash(img image.Image, hashLen uint) []byte {
	// gray scale
	img = grayscale(img)

	// resize
	img = resize.Resize(hashLen + 1, hashLen, img, resize.Lanczos3)

	diff := make([]byte, (hashLen+1) * hashLen)

	// create diff array
	for x := 0; x < img.Bounds().Size().X-1; x++ {
		for y := 0; y < img.Bounds().Size().Y; y++ {
			v1, _, _, _ := img.At(x, y).RGBA()
			v2, _, _, _ := img.At(x+1, y).RGBA()
			i := x + (img.Bounds().Size().X * y)
			if v1 < v2 {
				diff[i] = 1
			} else {
				diff[i] = 0
			}
		}
	}

	return diff
}

// generates a dhash for the given image and hash length.
// hash is based on the image data and not it's container.
// returns hash as hex string
func HashAsHex(img image.Image, hashLen uint) string {
	return hex.EncodeToString(Hash(img, hashLen))
}