package service

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo"
	"quizApp/model"
)

var fastjson jsoniter.API

func init() {
	fastjson = jsoniter.Config{
		EscapeHTML:             false,
		SortMapKeys:            false,
		ValidateJsonRawMessage: true,
	}.Froze()
}

// ApiResponse is used to respond with data and code
func ApiResponse(c echo.Context, code int, data interface{}) error {
	var resp struct {
		Error *model.ApiResponseError `json:"error"`
		Data  interface{}             `json:"data"`
	}

	if e, ok := data.(*model.ApiResponseError); ok {
		resp.Error = e
	} else if e, ok := data.(model.ApiResponseError); ok {
		resp.Error = &e
	} else {
		resp.Data = data
	}

	b, err := fastjson.Marshal(resp)
	if err != nil {
		return err
	}
	return c.JSONBlob(code, b)
}
