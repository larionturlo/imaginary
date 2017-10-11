package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestReadParamsFromJSON(t *testing.T) {
	q := `{"width": "100",
		"height": "80",
		"noreplicate": "1",
		"opacity": "0.2",
		"text": "hello",
		"background": "255,10,20",
		"operations":[
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
		]
	}`

	data := make(ParamsJSONScheme)

	if err := json.Unmarshal([]byte(q), &data); err != nil {
		fmt.Println(err)
	}

	params := readParamsFromJSON(data)

	assert := params.Width == 100 &&
		params.Height == 80 &&
		params.NoReplicate == true &&
		params.Opacity == 0.2 &&
		params.Text == "hello" &&
		params.Background[0] == 255 &&
		params.Background[1] == 10 &&
		params.Background[2] == 20 &&
		params.Operations[0].Name == "crop" &&
		params.Operations[1].Params["type"] == "webp"

	if assert == false {
		t.Error("Invalid params")
	}
}

func TestSetParamsRelatedImage(t *testing.T) {
	rawData := []byte(`{
		"width": "300",
		"height": "300"
		}`)
	data := make(ParamsJSONScheme)

	if err := json.Unmarshal(rawData, &data); err != nil {
		fmt.Println(err)
	}

	buf, _ := ioutil.ReadAll(readFile("imaginary.jpg"))

	params := SetParamsRelatedImage(buf, data)

	assert := params.Width == 0 &&
		params.Height == 300 &&
		params.NoCrop == true

	if assert == false {
		t.Errorf("Invalid params %d / %d", params.Width, params.Height)
	}
}
