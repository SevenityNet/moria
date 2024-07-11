package main

import (
	"bytes"
	"errors"
	"image"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
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

var videoEncodingQueue = make(chan *encodeVideoParams, 1000000)

type convertResult struct {
	OutputFileID string
	OutputData   []byte
	TmpFile      string
}

// Converts the given file to its compressed format, returning the extension and the compressed data.
func convertFile(in []byte, from fileType, wantedSubfolder, currExt string) (*convertResult, error) {
	if from == IMAGE {
		return convertImageFile(in, currExt)
	} else if from == VIDEO {
		return convertVideoFile(in, wantedSubfolder, currExt)
	} else if from == AUDIO {
		return convertAudioFile(in, currExt)
	}

	return nil, errors.New("invalid file type")
}

func convertImageFile(in []byte, currExt string) (*convertResult, error) {
	if currExt == "webp" {
		return &convertResult{
			OutputFileID: strings.ReplaceAll(uuid.NewString()+uuid.NewString(), "-", "") + ".webp",
			OutputData:   in,
			TmpFile:      "",
		}, nil
	}

	r := bytes.NewReader(in)
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	out := new(bytes.Buffer)
	if err := webpbin.Encode(out, img); err != nil {
		return nil, err
	}

	return &convertResult{
		OutputFileID: strings.ReplaceAll(uuid.NewString()+uuid.NewString(), "-", "") + ".webp",
		OutputData:   out.Bytes(),
		TmpFile:      "",
	}, nil
}

func convertVideoFile(in []byte, wantedSubfolder, currExt string) (*convertResult, error) {
	if currExt == "webm" {
		return &convertResult{
			OutputFileID: strings.ReplaceAll(uuid.NewString()+uuid.NewString(), "-", "") + ".webm",
			OutputData:   in,
			TmpFile:      "",
		}, nil
	}

	fileIdWoExt := strings.ReplaceAll(uuid.NewString()+uuid.NewString(), "-", "")
	inFileName := fileIdWoExt + "." + currExt

	if err := writeFileToDisk(wantedSubfolder, inFileName, in); err != nil {
		return nil, err
	}

	if currExt == "webm" {
		return &convertResult{
			OutputFileID: inFileName,
			OutputData:   nil,
			TmpFile:      "",
		}, nil
	}

	fileID := fileIdWoExt + ".webm"

	enqueueVideo(&encodeVideoParams{
		FileID:          fileID,
		WantedSubfolder: wantedSubfolder,
		InFileName:      inFileName,
	})

	return &convertResult{
		OutputFileID: fileID,
		OutputData:   nil,
		TmpFile:      "",
	}, nil
}

func convertAudioFile(in []byte, currExt string) (*convertResult, error) {
	if currExt == "aac" {
		return &convertResult{
			OutputFileID: strings.ReplaceAll(uuid.NewString()+uuid.NewString(), "-", "") + ".aac",
			OutputData:   in,
			TmpFile:      "",
		}, nil
	}

	inTmp, err := writeTmpFile(in, currExt)
	if err != nil {
		return nil, err
	}

	outTmp := getTmpFilePath("aac")

	if err := ffmpeg_go.Input(inTmp).Output(outTmp, ffmpeg_go.KwArgs{"c:a": "aac", "strict": "experimental"}).OverWriteOutput().Run(); err != nil {
		return nil, err
	}

	_ = deleteTmpFile(inTmp)

	return &convertResult{
		OutputFileID: strings.ReplaceAll(uuid.NewString()+uuid.NewString(), "-", "") + ".aac",
		OutputData:   nil,
		TmpFile:      outTmp,
	}, nil
}

func getFileType(file *multipart.FileHeader) (fileType, bool) {
	extension := strings.ToLower(filepath.Ext(file.Filename)[1:])

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

type encodeVideoParams struct {
	WantedSubfolder string
	InFileName      string
	FileID          string
}

func encodeVideo(params *encodeVideoParams) {
	outTmp := getTmpFilePath("webm")

	if err := ffmpeg_go.Input(getFilePath(params.WantedSubfolder, params.InFileName)).Output(outTmp, ffmpeg_go.KwArgs{"c:v": "libvpx-vp9", "b:v": "0", "crf": "30", "strict": "experimental"}).OverWriteOutput().Run(); err != nil {
		log.Println(err)
		return
	}

	if err := moveTmpToUploads(outTmp, params.WantedSubfolder, params.FileID); err != nil {
		log.Println(err)
		return
	}

	_ = deleteFileFromDisk(params.WantedSubfolder, params.InFileName)
}

func enqueueVideo(v *encodeVideoParams) {
	select {
	case videoEncodingQueue <- v:
	default:
		panic("video encoding queue full")
	}
}

func videoEncodingConsumer() {
	for {
		params := <-videoEncodingQueue
		encodeVideo(params)
	}
}
