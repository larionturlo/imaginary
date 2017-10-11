package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const SpecialOp = "double_resize_adaptive"

type ImageProcessConfig struct {
	ImageDir string
	ThumbDir string
}

type Task struct {
	ID        string           `json:"id"`
	SourceURL string           `json:"url"`
	Operation string           `json:"operation"`
	Params    ParamsJSONScheme `json:"params"`
}

type ResultProcessing struct {
	ID        string `json:"id"`
	Operation string `json:"operation"`
	URL       string `json:"url"`
	Error     string `json:"error"`
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

func RunProcess(taskData string) ResultProcessing {
	task, errTask := readTask(taskData)

	return ResultProcessing{task.ID, task.Operation, "new.example.com", errTask.Error()}
}

func RunImageProcess(sourceURL, operation string, params ParamsJSONScheme) error {

	buf, errRequest := requestImage("GET", sourceURL)
	if errRequest != nil {
		return errRequest
	}

	opts := readParamsFromJSON(params)

	return ImageProcess(buf, operation, opts, "./fixtures/")
}

func ImageProcess(buf []byte, operation string, opts ImageOptions, dir string) error {

	o := Operation(ImageOperations[operation])

	image, errProcess := o.Run(buf, opts)

	if errProcess != nil {
		return errProcess
	}
	return saveImageToFile(image, dir)
}

func saveImageToFile(img Image, imgDir string) error {
	im, _, errDecode := image.Decode(bytes.NewReader(img.Body))
	if errDecode != nil {
		return fmt.Errorf("Error image decode: %v", errDecode)
	}
	hash := fmt.Sprintf("%x", img.Hash)
	dir := imgDir + hash[:2]

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0775)
	}

	toimg, errCreateFile := os.Create(dir + "/" + hash + ".jpg")
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

func readTask(taskData string) (Task, error) {
	task := Task{}
	error := json.Unmarshal([]byte(taskData), &task)
	return task, error
}
