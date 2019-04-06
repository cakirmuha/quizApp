package main

import (
	"net/http"
	"quizApp/cmd/quiz/assets"
	"quizApp/service"

	"github.com/labstack/echo"
)

func (s *ServerContext) SetupHandlers() {
	const handlerPrefix = "/api/v1"

	api, endSetupFunc := s.Context.SetupHandlers(handlerPrefix)

	var (
		jsonBody = service.EnsureContentTypeFunc(echo.MIMEApplicationJSON)
	)

	{
		fs := http.FileServer(assets.Assets)
		api.GET("/assets/*", echo.WrapHandler(http.StripPrefix(handlerPrefix+"/assets/", fs)), service.NoCacheMiddleware).Name = "GetStaticAssets"
	}
	{
		g := api.Group("", jsonBody)
		g.POST("/user", addUser).Name = "addUser"
		g.POST("/user/:username/answer", answerQuestions).Name = "answerQuestions"
	}
	{
		g := api.Group("")
		g.GET("/question", getQuestions).Name = "getQuestions"
	}
	endSetupFunc()
}
