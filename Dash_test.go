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

	if ajHash != "616e35d28e85921898" {
		t.Errorf("Expected aqua jpg hash to be 616e35d28e85921898 got %s", ajHash)
	}

	if apHash != "616e35d28e85921898" {
		t.Errorf("Expected aqua png hash to be 616e35d28e85921898 got %s", apHash)
	}

	if mpHash != "3731bcce426130c060" {
		t.Errorf("Expected megumin hash to be 3731bcce426130c060 got %s", mpHash)
	}

	if ajHash != apHash {
		t.Errorf("Expected both aqua hashes to be equal!\nGot: %s != %s", ajHash, apHash)
		return
	}

	if mpHash == ajHash {
		t.Errorf("Expected megumin hash to be different to aqua hash!\nGot %s == %s", mpHash, ajHash)
		return
	}

	dist, err := HammingDistanceHex(apHash, ajHash)
	if err != nil {
		t.Error(err)
	}

	if dist != 0 {
		t.Errorf("Expected distance between aqua images to be 0, got %d", dist)
	}

	dist, err = HammingDistanceHex(mpHash, ajHash)
	if err != nil {
		t.Error(err)
	}

	if dist == 0 {
		t.Errorf("Expected distance between aqua and megumin images to not be 0, got %d", dist)
	}

	meguBright, err := loadImage("test_imgs/megu_bright.png", "png")
	dist, err = HammingDistanceHex(mpHash, HashAsHex(meguBright, l))
	if err != nil {
		t.Error(err)
	}

	if dist != 2 {
		t.Errorf("Expected distance between megumin and megu_bright images to be 2, got %d", dist)
	}
}
