package main

import (
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

	params := readParamsFromJSON(q)

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
