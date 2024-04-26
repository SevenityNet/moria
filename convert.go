package main

import (
	"bytes"
	"errors"
	"image"
	"mime/multipart"
	"strings"

	"github.com/nickalie/go-webpbin"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

type fileType int

const (
	IMAGE fileType = iota
	VIDEO
	AUDIO
)

var (
	ALLOWED_IMAGE_EXTENSIONS = map[string]bool{
		"jpg":  true,
		"jpeg": true,
		"png":  true,
		"webp": true,
		"bmp":  true,
	}
	ALLOWED_VIDEO_EXTENSIONS = map[string]bool{
		"mp4":  true,
		"mkv":  true,
		"webm": true,
	}
	ALLOWED_AUDIO_EXTENSIONS = map[string]bool{
		"mp3": true,
		"ogg": true,
		"wav": true,
	}
)

type convertResult struct {
	OutputData []byte
	TmpFile    string
}

// Converts the given file to its compressed format, returning the extension and the compressed data.
func convertFile(in []byte, from fileType, currExt string) (string, *convertResult, error) {
	if from == IMAGE {
		return convertImageFile(in)
	} else if from == VIDEO {
		return convertVideoFile(in, currExt)
	} else if from == AUDIO {
		return convertAudioFile(in, currExt)
	}

	return "", nil, errors.New("invalid file type")
}

func convertImageFile(in []byte) (string, *convertResult, error) {
	r := bytes.NewReader(in)
	img, _, err := image.Decode(r)
	if err != nil {
		return "", nil, err
	}

	out := new(bytes.Buffer)
	if err := webpbin.Encode(out, img); err != nil {
		return "", nil, err
	}

	return "webp", &convertResult{
		OutputData: out.Bytes(),
		TmpFile:    "",
	}, nil
}

func convertVideoFile(in []byte, currExt string) (string, *convertResult, error) {
	inTmp, err := writeTmpFile(in, currExt)
	if err != nil {
		return "", nil, err
	}

	outTmp := getTmpFilePath("webm")

	if err := ffmpeg_go.Input(inTmp).Output(outTmp, ffmpeg_go.KwArgs{"c:v": "libvpx-vp9", "b:v": "0", "crf": "30", "strict": "experimental"}).OverWriteOutput().Run(); err != nil {
		return "", nil, err
	}

	_ = deleteTmpFile(inTmp)

	return "webm", &convertResult{
		OutputData: nil,
		TmpFile:    outTmp,
	}, nil
}

func convertAudioFile(in []byte, currExt string) (string, *convertResult, error) {
	inTmp, err := writeTmpFile(in, currExt)
	if err != nil {
		return "", nil, err
	}

	outTmp := getTmpFilePath("aac")

	if err := ffmpeg_go.Input(inTmp).Output(outTmp, ffmpeg_go.KwArgs{"c:a": "aac", "strict": "experimental"}).OverWriteOutput().Run(); err != nil {
		return "", nil, err
	}

	_ = deleteTmpFile(inTmp)

	return "aac", &convertResult{
		OutputData: nil,
		TmpFile:    outTmp,
	}, nil
}

func getFileType(file *multipart.FileHeader) (fileType, bool) {
	extension := strings.Split(file.Filename, ".")[1]

	if ALLOWED_IMAGE_EXTENSIONS[extension] {
		return IMAGE, true
	}

	if ALLOWED_VIDEO_EXTENSIONS[extension] {
		return VIDEO, true
	}

	if ALLOWED_AUDIO_EXTENSIONS[extension] {
		return AUDIO, true
	}

	return -1, false
}

func getFileTypeByFilePath(filePath string) (fileType, bool) {
	if strings.HasSuffix(filePath, ".webp") {
		return IMAGE, true
	}

	if strings.HasSuffix(filePath, ".webm") {
		return VIDEO, true
	}

	if strings.HasSuffix(filePath, ".aac") {
		return AUDIO, true
	}

	return -1, false
}
