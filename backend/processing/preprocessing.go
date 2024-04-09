package processing

import (
	"moria/config"

	"github.com/h2non/bimg"
)

// Applies compression to the given image by converting it to WebP format.
func Compress(in []byte) ([]byte, error) {
	if !config.IsProcessingCompressionEnabled() {
		return in, nil
	}

	img := bimg.NewImage(in)
	newImg, err := img.Convert(bimg.WEBP)
	if err != nil {
		return nil, err
	}

	return newImg, nil
}
