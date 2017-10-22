package main

import (
	"crypto/md5"
	"fmt"
	"reflect"
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

	task, error := readTask([]byte(taskData))
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

	img := Image{Body: buf, Mime: "image/jpeg"}
	hash := fmt.Sprintf("%x", md5.Sum(buf))
	errSave := saveImageToFile(img, "./fixtures/", hash)
	if errSave != nil {
		t.Errorf("Save error: %v", errSave)
	}
}

func TestRunProcess(t *testing.T) {
	taskDataGood := `{
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
	taskGood, _ := readTask([]byte(taskDataGood))

	type args struct {
		taskData Task
	}
	tests := []struct {
		name  string
		args  args
		want  ImageResultQueueMSG
		want1 error
	}{
		{
			name: "taskDataGood",
			args: args{taskGood},
			want: ImageResultQueueMSG{
				ID:        taskGood.ID,
				Operation: taskGood.Operation,
			},
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := RunProcess(tt.args.taskData)

			if reflect.ValueOf(got).Type() != reflect.ValueOf(tt.want).Type() {
				t.Errorf("RunProcess() got1 = %v, want %v", reflect.ValueOf(got1).Type(), reflect.ValueOf(tt.want1).Type())
			}
			if got.ID != tt.want.ID {
				t.Errorf("RunProcess() got1 = %v, want %v", got.ID, tt.want.ID)
			}
			if got.Operation != tt.want.Operation {
				t.Errorf("RunProcess() got1 = %v, want %v", got.Operation, tt.want.Operation)
			}

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("RunProcess() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_requestImage(t *testing.T) {
	url := `https://image.63pokupki.ru/images/0f/0f020c12f8825.jpg`
	method := "GET"
	type args struct {
		method    string
		sourceURL string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "loadImg",
			args:    args{method, url},
			want:    []byte("123"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := requestImage(tt.args.method, tt.args.sourceURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("requestImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.ValueOf(got).Type() != reflect.ValueOf(tt.want).Type() {
				t.Errorf("requestImage() = %v, want %v", got, tt.want)
			}
		})
	}
}
