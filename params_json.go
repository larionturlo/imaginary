package main

import (
	"encoding/json"
	"strings"

	bimg "gopkg.in/h2non/bimg.v1"
)

// not map[string]interface{} , because param operations return []interface{}, not string
type ParamsJSONScheme map[string]*json.RawMessage

func readParamsFromJSON(data ParamsJSONScheme) ImageOptions {
	params := make(map[string]interface{})

	for key, kind := range allowedParams {
		param := findParamInMapJSON(key, data)
		params[key] = parseParam(param, kind)
	}

	return mapImageParams(params)
}

func SetParamsRelatedImage(buf []byte, data ParamsJSONScheme) ImageOptions {
	params := make(map[string]interface{})

	meta, err := bimg.Metadata(buf)
	if err != nil {
		return mapImageParams(params)
	}

	for key, kind := range allowedParams {
		params[key] = parseParam("", kind)
	}

	var param string
	// todo implement checking exists width and height
	if meta.Size.Width > meta.Size.Height {
		json.Unmarshal(*data["width"], &param)
		params["width"] = parseParam(param, "int")
	} else {
		json.Unmarshal(*data["height"], &param)
		params["height"] = parseParam(param, "int")
	}
	params["nocrop"] = true

	return mapImageParams(params)
}

func findParamInMapJSON(key string, mapJSON ParamsJSONScheme) string {
	if foundParam := mapJSON[key]; foundParam != nil {
		jsonParam, error := foundParam.MarshalJSON() // because it`s only way for correctly convert array json.RawMessage to string
		if error != nil {
			return ""
		}
		return strings.Trim(string(jsonParam), `"`)
	}
	return ""

}
