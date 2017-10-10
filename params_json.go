package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// not map[string]interface{} , because param operations return []interface{}, not string
type ParamsJSONScheme map[string]*json.RawMessage

func readParamsFromJSON(rawData string) ImageOptions {
	params := make(map[string]interface{})
	data := make(ParamsJSONScheme)

	if error := json.Unmarshal([]byte(rawData), &data); error != nil {
		fmt.Println(error)
	}

	for key, kind := range allowedParams {
		param := findParamInMapJSON(key, data)
		params[key] = parseParam(param, kind)
	}

	return mapImageParams(params)
}

func findParamInMapJSON(key string, mapJSON ParamsJSONScheme) string {
	if foundParam := mapJSON[key]; foundParam != nil {
		jsonParam, error := foundParam.MarshalJSON() // because it`s only way for correctly convert json.RawMessage to string
		if error != nil {
			return ""
		}
		return strings.Trim(string(jsonParam), `"`)
	}
	return ""

}
