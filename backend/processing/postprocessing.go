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
	options := bimg.Options{}
	update := false

	// Resize
	if config.IsProcessingResizeEnabled() {
		w := getQueryInt(c, "w", nil)
		h := getQueryInt(c, "h", nil)
		keepAspectRatio := getQueryBool(c, "keepAspectRatio", _p(false))
		if w != nil && h != nil {
			options.Width = *w
			options.Height = *h
			if *keepAspectRatio {
				options.Embed = true
			} else {
				options.Force = true
			}
			update = true
		}
	}

	// Rotate
	if config.IsProcessingRotateEnabled() {
		degree := getQueryInt(c, "rotate", nil)
		if degree != nil {
			options.Rotate = bimg.Angle(*degree)
			update = true
		}
	}

	// Crop
	if config.IsProcessingCropEnabled() {
		w := getQueryInt(c, "cw", nil)
		h := getQueryInt(c, "ch", nil)
		gravity := getQueryInt(c, "gravity", nil)
		if w != nil && h != nil && gravity != nil {
			options.Width = *w
			options.Height = *h
			options.Gravity = bimg.Gravity(*gravity)
			options.Crop = true
			update = true
		}
	}

	// To Grayscale
	if config.IsProcessingToGrayscaleEnabled() {
		toGrayscale := getQueryBool(c, "bw", _p(false))
		if *toGrayscale {
			options.Interpretation = bimg.InterpretationBW
			update = true
		}
	}

	// Watermark currently error producing
	/*if config.IsProcessingWatermarkEnabled() {
		text := c.Query("watermarkText")
		if text != "" {
			font := c.Query("watermarkFont")
			if font == "" {
				font = "sans bold 12"
			}

			metaData, err := img.Metadata()
			if err != nil {
				return nil, err
			}

			options.Watermark = bimg.Watermark{
				Text:    text,
				Opacity: *getQueryFloat32(c, "watermarkOpacity", _p[float32](1.0)),
				Width:   *getQueryInt(c, "watermarkWidth", _p(metaData.Size.Width)),
				DPI:     *getQueryInt(c, "watermarkDPI", _p(72)),
				Margin:  *getQueryInt(c, "watermarkMargin", _p(10)),
				Font:    font,
			}
			update = true
		}
	}*/

	// Flip
	if config.IsProcessingFlipEnabled() {
		flip := getQueryBool(c, "flip", _p(false))
		if *flip {
			options.Flip = true
			update = true
		}
	}

	// Flop
	if config.IsProcessingFlopEnabled() {
		flop := getQueryBool(c, "flop", _p(false))
		if *flop {
			options.Flop = true
			update = true
		}
	}

	// Zoom
	if config.IsProcessingZoomEnabled() {
		factor := getQueryInt(c, "zoom", nil)
		if factor != nil {
			options.Zoom = *factor
			update = true
		}
	}

	// Blur
	if config.IsProcessingBlurEnabled() {
		sigma := getQueryFloat64(c, "blurSigma", _p(1.0))
		minAmpl := getQueryFloat64(c, "blurMinAmpl", _p(1.0))
		if sigma != nil && minAmpl != nil {
			options.GaussianBlur = bimg.GaussianBlur{
				Sigma:   *sigma,
				MinAmpl: *minAmpl,
			}
			update = true
		}
	}

	if !update {
		return in, nil
	}

	processedImage, err := img.Process(options)
	if err != nil {
		return nil, err
	}

	return processedImage, nil
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
