package processing

import (
	"moria/config"

	"github.com/h2non/bimg"
)

// Applies compression to the given image by converting it to WebP format.
func Compress(in []byte) ([]byte, bool, error) {
	if !config.IsProcessingCompressionEnabled() {
		return in, false, nil
	}

	img := bimg.NewImage(in)
	newImg, err := img.Convert(bimg.WEBP)
	if err != nil {
		return nil, false, err
	}

	return newImg, true, nil
}
