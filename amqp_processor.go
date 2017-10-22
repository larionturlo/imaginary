package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type ImageProcessConfig struct {
	ImageLocalRoot string
	ImageDir       string
	ThumbDir       string
	ImageHost      string
}

type ImageResultQueueMSG struct {
	ID        string `json:"id"`
	Operation string `json:"operation"`
	Hash      string `json:"hash"`
}

type Task struct {
	ID        string           `json:"id"`
	SourceURL string           `json:"url"`
	Operation string           `json:"operation"`
	Params    ParamsJSONScheme `json:"params"`
}

var ImageOperations = map[string]Operation{
	"resize":    Resize,
	"enlarge":   Enlarge,
	"extract":   Extract,
	"crop":      Crop,
	"smartcrop": SmartCrop,
	"rotate":    Rotate,
	"flip":      Flip,
	"flop":      Flop,
	"thumbnail": Thumbnail,
	"zoom":      Zoom,
	"convert":   Convert,
	"watermark": Watermark,
	"info":      Info,
	"blur":      GaussianBlur,
	"pipeline":  Pipeline,
}

func RunProcess(task Task) (ImageResultQueueMSG, error) {
	settings := ImageProcessConfig{
		ImageLocalRoot: "./fixtures/",
		ImageDir:       "image/",
		ThumbDir:       "thumb/",
		ImageHost:      "image.63pokupki.ru",
	}

	buf, errRequest := requestImage("GET", task.SourceURL)
	if errRequest != nil {
		return ImageResultQueueMSG{}, errRequest
	}

	opts := SetParamsRelatedImage(buf, task.Params)
	hash := fmt.Sprintf("%x", md5.Sum(buf))
	dir := settings.ImageLocalRoot + settings.ImageDir + hash[:2]

	imageProcess(buf, task.Operation, opts, dir, hash)

	thumbDir := settings.ImageLocalRoot + settings.ThumbDir + hash[:2]
	thumbParams := make(ParamsJSONScheme)
	thumbParams["width"] = task.Params["thumb_width"]
	thumbParams["height"] = task.Params["thumb_height"]
	thumbOpts := SetParamsRelatedImage(buf, thumbParams)

	imageProcess(buf, task.Operation, thumbOpts, thumbDir, hash)

	return ImageResultQueueMSG{task.ID, task.Operation, hash}, nil
}

func imageProcess(buf []byte, operation string, opts ImageOptions, dir, hash string) error {

	o := Operation(ImageOperations[operation])

	image, errProcess := o.Run(buf, opts)

	if errProcess != nil {
		return errProcess
	}
	return saveImageToFile(image, dir, hash)
}

func saveImageToFile(img Image, imgDir, hash string) error {
	im, _, errDecode := image.Decode(bytes.NewReader(img.Body))
	if errDecode != nil {
		return fmt.Errorf("Error image decode: %v", errDecode)
	}

	if _, err := os.Stat(imgDir); os.IsNotExist(err) {
		os.Mkdir(imgDir, 0775)
	}

	toimg, errCreateFile := os.Create(imgDir + "/" + hash + ".jpg")
	if errCreateFile != nil {
		return fmt.Errorf("Error create new file: %v", errCreateFile)
	}
	defer toimg.Close()
	err := jpeg.Encode(toimg, im, &jpeg.Options{jpeg.DefaultQuality})
	if err != nil {
		return fmt.Errorf("Error image jpeg encode: %v", err)
	}
	return nil
}

func requestImage(method, sourceURL string) ([]byte, error) {
	url, errParse := url.Parse(sourceURL)
	if errParse != nil {
		return nil, fmt.Errorf("Error parse url: %v", errParse)
	}

	req, _ := http.NewRequest(method, url.String(), nil)
	req.Header.Set("User-Agent", "imaginary/"+Version)
	req.URL = url
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error downloading image: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Error downloading image: (status=%d) (url=%s)", res.StatusCode, req.URL.String())
	}

	// Read the body
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Unable to create image from response body: %s (url=%s)", req.URL.String(), err)
	}
	return buf, nil
}

func readTask(taskData []byte) (Task, error) {
	task := Task{}
	error := json.Unmarshal(taskData, &task)
	return task, error
}
