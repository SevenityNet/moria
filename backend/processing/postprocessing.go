package processing

import (
	"moria/config"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/h2non/bimg"
)

func PostProcess(in []byte, c *gin.Context) ([]byte, error) {
	if !config.IsProcessingEnabled() {
		return in, nil
	}

	img := bimg.NewImage(in)

	if config.IsProcessingResizeEnabled() {
		w := getQueryInt(c, "w", nil)
		h := getQueryInt(c, "h", nil)
		keepAspectRatio := getQueryBool(c, "keepAspectRatio", _p(false))

		if w != nil && h != nil {
			i, err := resizeImage(img, *w, *h, *keepAspectRatio)
			if err != nil {
				return nil, err
			}

			img = i
		}
	}

	if config.IsProcessingRotateEnabled() {
		degree := getQueryInt(c, "rotate", nil)
		if degree != nil {
			i, err := rotateImage(img, bimg.Angle(*degree))
			if err != nil {
				return nil, err
			}

			img = i
		}
	}

	if config.IsProcessingCropEnabled() {
		w := getQueryInt(c, "cw", nil)
		h := getQueryInt(c, "ch", nil)
		gravity := getQueryInt(c, "gravity", nil)
		if w != nil && h != nil && gravity != nil {
			i, err := cropImage(img, *w, *h, bimg.Gravity(*gravity))
			if err != nil {
				return nil, err
			}

			img = i
		}
	}

	if config.IsProcessingToGrayscaleEnabled() {
		toGrayscale := getQueryBool(c, "bw", _p(false))
		if *toGrayscale {
			i, err := imageToGrayscale(img)
			if err != nil {
				return nil, err
			}

			img = i
		}
	}

	if config.IsProcessingBlurEnabled() {
		sigma := getQueryFloat64(c, "blurSigma", nil)
		minAmpl := getQueryFloat64(c, "blurMinAmpl", nil)
		if sigma != nil && minAmpl != nil {
			i, err := blurImage(img, *sigma, *minAmpl)
			if err != nil {
				return nil, err
			}

			img = i
		}
	}

	if config.IsProcessingWatermarkEnabled() {
		font := c.Query("watermarkFont")
		if font == "" {
			font = "sans bold 12"
		}

		watermark := bimg.Watermark{
			Text:    c.Query("watermarkText"),
			Opacity: *getQueryFloat32(c, "watermarkOpacity", _p[float32](1.0)),
			Width:   *getQueryInt(c, "watermarkWidth", _p(100)),
			DPI:     *getQueryInt(c, "watermarkDPI", _p(72)),
			Margin:  *getQueryInt(c, "watermarkMargin", _p(10)),
			Font:    font,
		}
		i, err := watermarkImage(img, watermark)
		if err != nil {
			return nil, err
		}

		img = i
	}

	if config.IsProcessingFlipEnabled() {
		flip := getQueryBool(c, "flip", _p(false))
		if *flip {
			i, err := flipImage(img)
			if err != nil {
				return nil, err
			}

			img = i
		}
	}

	if config.IsProcessingFlopEnabled() {
		flip := getQueryBool(c, "flop", _p(false))
		if *flip {
			i, err := flopImage(img)
			if err != nil {
				return nil, err
			}

			img = i
		}
	}

	if config.IsProcessingZoomEnabled() {
		factor := getQueryInt(c, "zoom", nil)
		if factor != nil {
			i, err := zoomImage(img, *factor)
			if err != nil {
				return nil, err
			}

			img = i
		}
	}

	return img.Image(), nil
}

func resizeImage(in *bimg.Image, width, height int, keepAspectRatio bool) (*bimg.Image, error) {
	if keepAspectRatio {
		out, err := in.Resize(width, height)
		if err != nil {
			return nil, err
		}

		return bimg.NewImage(out), nil
	} else {
		out, err := in.ForceResize(width, height)
		if err != nil {
			return nil, err
		}

		return bimg.NewImage(out), nil
	}
}

func rotateImage(in *bimg.Image, degree bimg.Angle) (*bimg.Image, error) {
	cc := degree < 0
	degree = degree * -1

	out, err := in.Rotate(degree)
	if err != nil {
		return nil, err
	}

	if cc {
		out, err = bimg.NewImage(out).Flip()
		if err != nil {
			return nil, err
		}
	}

	return bimg.NewImage(out), nil
}

func cropImage(in *bimg.Image, width, height int, gravity bimg.Gravity) (*bimg.Image, error) {
	out, err := in.Crop(width, height, gravity)
	if err != nil {
		return nil, err
	}

	return bimg.NewImage(out), nil
}

func imageToGrayscale(in *bimg.Image) (*bimg.Image, error) {
	out, err := in.Colourspace(bimg.InterpretationBW)
	if err != nil {
		return nil, err
	}

	return bimg.NewImage(out), nil
}

func blurImage(in *bimg.Image, sigma, minAmpl float64) (*bimg.Image, error) {
	out, err := in.Process(bimg.Options{
		GaussianBlur: bimg.GaussianBlur{
			Sigma:   sigma,
			MinAmpl: minAmpl,
		},
	})
	if err != nil {
		return nil, err
	}

	return bimg.NewImage(out), nil
}

func watermarkImage(in *bimg.Image, watermark bimg.Watermark) (*bimg.Image, error) {
	out, err := in.Watermark(watermark)
	if err != nil {
		return nil, err
	}

	return bimg.NewImage(out), nil
}

func flipImage(in *bimg.Image) (*bimg.Image, error) {
	out, err := in.Flip()
	if err != nil {
		return nil, err
	}

	return bimg.NewImage(out), nil
}

func flopImage(in *bimg.Image) (*bimg.Image, error) {
	out, err := in.Flop()
	if err != nil {
		return nil, err
	}

	return bimg.NewImage(out), nil
}

func zoomImage(in *bimg.Image, factor int) (*bimg.Image, error) {
	out, err := in.Zoom(factor)
	if err != nil {
		return nil, err
	}

	return bimg.NewImage(out), nil
}

func getQueryInt(c *gin.Context, key string, defaultVal *int) *int {
	value, exists := c.GetQuery(key)
	if !exists {
		return defaultVal
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultVal
	}

	return &intValue
}

func getQueryFloat64(c *gin.Context, key string, defaultVal *float64) *float64 {
	value, exists := c.GetQuery(key)
	if !exists {
		return defaultVal
	}

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultVal
	}

	return &floatValue
}

func getQueryFloat32(c *gin.Context, key string, defaultVal *float32) *float32 {
	value, exists := c.GetQuery(key)
	if !exists {
		return defaultVal
	}

	floatValue, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return defaultVal
	}

	return _p(float32(floatValue))
}

func getQueryBool(c *gin.Context, key string, defaultVal *bool) *bool {
	value, exists := c.GetQuery(key)
	if !exists {
		if c.Request.URL.Query().Has(key) {
			return _p(true)
		}

		return defaultVal
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultVal
	}

	return &boolValue
}

func _p[T any](v T) *T {
	return &v
}
