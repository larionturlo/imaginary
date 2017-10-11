package main

import (
	"crypto/md5"
	"testing"
)

func TestReadTask(t *testing.T) {
	operationsData := `[
		{
			"operation": "crop",
			"params": {
			"width": 500,
			"height": 300
			}
		},
		{
			"operation": "convert",
			"params": {
			"type": "webp"
			}
		}
	]`
	taskData := `{
		"id":"000000",
		"url":"example.com/uri",
		"operation": "pipeline",
		"params":{
			"width": "100",
			"height": "80",
			"operations":` + operationsData + `
		}
	}`

	task, error := readTask(taskData)
	if error != nil {
		t.Error(error)
	}

	assert := task.ID == "000000" &&
		task.SourceURL == "example.com/uri" &&
		task.Operation == "pipeline"

	width, error := task.Params["width"].MarshalJSON()
	if error != nil {
		t.Error(error)
	}

	assert = assert && string(width) == `"100"`

	operations, error := task.Params["operations"].MarshalJSON()
	if error != nil {
		t.Error(error)
	}

	assert = assert && string(operations) == operationsData

	if assert == false {
		t.Error("Invalid params")
	}
}

func TestRequestImageAndSave(t *testing.T) {
	url := `https://image.63pokupki.ru/images/0f/0f020c12f8825.jpg`
	method := "GET"
	buf, errRequest := requestImage(method, url)
	if errRequest != nil {
		t.Errorf("Request error: %v", errRequest)
	}

	img := Image{Body: buf, Mime: "image/jpeg", Hash: md5.Sum(buf)}
	errSave := saveImageToFile(img, "./fixtures/")
	if errSave != nil {
		t.Errorf("Save error: %v", errSave)
	}
}

func TestRunImageProcess(t *testing.T) {
	taskData := `{
		"id":"000000",
		"url":"https://image.63pokupki.ru/images/0f/0f020c12f8825.jpg",
		"operation": "smartcrop",
		"params":{
			"width": "100",
			"height": "100",
			"nocrop": "1"
		}
	}`

	task, errReadTask := readTask(taskData)
	if errReadTask != nil {
		t.Error(errReadTask)
	}

	err := RunImageProcess(task.SourceURL, task.Operation, task.Params)
	if err != nil {
		t.Error(err)
	}
}

func TestRunProcess(t *testing.T) {
	taskData := `{
		"id":"000000",
		"url":"https://image.63pokupki.ru/images/0f/0f020c12f8825.jpg",
		"operation": "resize",
		"params":{
			"width": "300",
			"height": "300",
			"thumb_width": "150",
			"thumb_height": "150"
		}
	}`

	RunProcess(taskData)
}
