package dash

import (
	"encoding/hex"
	"errors"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"strconv"
)

type bitString string

func (b bitString) asByteSlice() []byte {
	var out []byte
	var str string

	for i := len(b); i > 0; i -= 8 {
		if i-8 < 0 {
			str = string(b[0:i])
		} else {
			str = string(b[i-8 : i])
		}
		v, err := strconv.ParseUint(str, 2, 8)
		if err != nil {
			panic(err)
		}
		out = append([]byte{byte(v)}, out...)
	}
	return out
}

// converts image to grayscale
func grayscale(img image.Image) image.Image {
	gray := image.NewGray(img.Bounds())
	for x := 0; x < img.Bounds().Size().X; x++ {
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
	img = resize.Resize(hashLen+1, hashLen, img, resize.Lanczos3)

	diffBits := make([]int8, (hashLen+1)*hashLen)

	// create diff array
	for x := 0; x < img.Bounds().Size().X-1; x++ {
		for y := 0; y < img.Bounds().Size().Y; y++ {
			v1, _, _, _ := img.At(x, y).RGBA()
			v2, _, _, _ := img.At(x+1, y).RGBA()
			i := x + (img.Bounds().Size().X * y)
			if v1 < v2 {
				diffBits[i] = 1
			} else {
				diffBits[i] = 0
			}
		}
	}

	diffStr := ""
	for _, bit := range diffBits {
		diffStr += strconv.Itoa(int(bit))
	}

	bytes := bitString(diffStr).asByteSlice()
	return bytes
}

// generates a dhash for the given image and hash length.
// hash is based on the image data and not it's container.
// returns hash as hex string
func HashAsHex(img image.Image, hashLen uint) string {
	return hex.EncodeToString(Hash(img, hashLen))
}

func HammingDistance(bytesA []byte, bytesB []byte) (int, error) {
	if len(bytesA) != len(bytesB) {
		return -1, errors.New("bytes A and B are of differing lengths")
	}

	diff := 0
	for i := 0; i < len(bytesA); i++ {
		b1 := bytesA[i]
		b2 := bytesB[i]
		// use byte mask to compare individual bits
		for j := 0; j < 8; j++ {
			mask := byte(1 << uint(j))
			if (b1 & mask) != (b2 & mask) {
				diff += 1
			}
		}
	}

	return diff, nil
}

func HammingDistanceHex(a string, b string) (int, error) {
	return HammingDistance([]byte(a), []byte(b))
}
