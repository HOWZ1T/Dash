package dash

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"testing"
)

func loadImage(path string, imgType string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var decode_fn func(r io.Reader) (image.Image, error)
	if imgType == "png" {
		decode_fn = png.Decode
	} else if imgType == "jpg" || imgType == "jpeg" {
		decode_fn = jpeg.Decode
	} else {
		return nil, fmt.Errorf("unhandled type: %s", imgType)
	}

	img, err := decode_fn(f)
	if err != nil {
		return nil, err
	}

	err = f.Close()
	if err != nil {
		return nil, err
	}

	return img, nil
}

func TestHash(t *testing.T) {
	// load test images
	aquaJpg, err := loadImage("test_imgs/aqua.jpg", "jpg")
	if err != nil {
		t.Error(err)
		return
	}

	aquaPng, err := loadImage("test_imgs/aqua.png", "png")
	if err != nil {
		t.Error(err)
		return
	}

	meguPng, err := loadImage("test_imgs/megumin.png", "png")
	if err != nil {
		t.Error(err)
		return
	}

	// generate hashes
	var l uint = 8
	ajHash := HashAsHex(aquaJpg, l)
	apHash := HashAsHex(aquaPng, l)
	mpHash := HashAsHex(meguPng, l)

	if ajHash != apHash {
		t.Errorf("Expected both aqua hashes to be equal!\nGot: %s != %s", ajHash, apHash)
		return
	}

	if mpHash == ajHash {
		t.Errorf("Expected megumin hash to be different to aqua hash!\nGot %s == %s", mpHash, ajHash)
		return
	}
}