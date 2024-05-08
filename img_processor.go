package main

import (
	"bytes"
	"image/color"
	"image/png"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/nickalie/go-webpbin"
)

func postprocessImage(c *gin.Context, inBytes []byte) ([]byte, error) {
	rotateQuery := c.Query("r")
	if rotateQuery != "" {
		degrees, err := strconv.Atoi(rotateQuery)
		if err != nil {
			return nil, err
		}

		inBytes, err = rotateImage(inBytes, degrees)
		if err != nil {
			return nil, err
		}
	}

	return inBytes, nil
}

// rotateImage rotates an image by the given degrees. If the degrees are less then 0, the image will be rotated counter-clockwise.
func rotateImage(inBytes []byte, degrees int) ([]byte, error) {
	webpDecoder := webpbin.NewDWebP()
	webpEncoder := webpbin.NewCWebP()

	decodedImg, err := webpDecoder.Input(bytes.NewReader(inBytes)).Run()
	if err != nil {
		return nil, err
	}

	dstImg := imaging.Rotate(decodedImg, float64(degrees), color.Transparent)

	var encodeBuf bytes.Buffer
	png.Encode(&encodeBuf, dstImg)

	var outBuf bytes.Buffer
	err = webpEncoder.Input(&encodeBuf).Output(&outBuf).Run()

	return outBuf.Bytes(), err
}
