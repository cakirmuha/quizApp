package main

import (
	"net/http"
	"quizApp/model"
	"quizApp/service"

	"github.com/labstack/echo"
)

func addUser(c echo.Context) error {
	cc := c.(*ServerContext)

	r := model.User{}
	if err := c.Bind(&r); err != nil {
		return model.NewApiError(model.ApiErrorBadRequest, "bind", err)
	}

	if err := cc.DB().AddUser(r); err != nil {
		return model.NewApiError(model.ApiErrorExists, err.Error(), nil)
	}

	return service.ApiResponse(c, http.StatusCreated, nil)

}
