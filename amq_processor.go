package main

import (
	"encoding/json"
)

type Task struct {
	ID        string                      `json:"id"`
	Operation string                      `json:"operation"`
	Params    map[string]*json.RawMessage `json:"params"`
}

type ResultProcessing struct {
	ID        string `json:"id"`
	Operation string `json:"operation"`
	URL       string `json:"url"`
}

const ImageOperations = map[string]Operation{
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

func RunProcessor(task Task) ResultProcessing {
	return ResultProcessing{"000000", task.Operation, "example.com"}
}
