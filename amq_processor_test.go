package main

import (
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
