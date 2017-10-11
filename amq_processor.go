package main

import "encoding/json"

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

// func RunProcess(taskData string) ResultProcessing {
// 	task := readTask(taskData)
// 	return ResultProcessing{task.ID, task.Operation, "new.example.com"}
// }

// func RunImageProcess(operation string, params ParamsJSONScheme) {

// 	opts := readParamsFromJSON(params)

// 	o := Operation(ImageOperations[operation])

// 	imgSource = NewHttpImageSource()

// 	image, err := o.Run(buf, opts)
// 	if err != nil {
// 		ErrorReply(r, w, NewError("Error while processing the image: "+err.Error(), BadRequest), o)
// 		return
// 	}
// }
func readTask(taskData string) (Task, error) {
	task := Task{}
	error := json.Unmarshal([]byte(taskData), &task)
	return task, error
}
