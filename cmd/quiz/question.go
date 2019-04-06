package main

import (
	"fmt"
	"net/http"
	"quizApp/model"
	"quizApp/service"

	"github.com/labstack/echo"
)

func getQuestions(c echo.Context) error {
	cc := c.(*ServerContext)

	ques := cc.DB().GetQuestions()
	if len(ques) == 0 {
		return model.NewApiError(model.ApiErrorNotFound, "No question found", fmt.Errorf("No question found"))
	}

	return service.ApiResponse(c, http.StatusOK, ques)
}

func answerQuestions(c echo.Context) error {
	cc := c.(*ServerContext)

	if ok := cc.DB().GetUser(c.Param("username")); !ok {
		return model.NewApiError(model.ApiErrorNotFound, "User not found", nil)
	}

	var r []model.Option

	if err := c.Bind(&r); err != nil {
		return model.NewApiError(model.ApiErrorBadRequest, "bind", err)
	}

	ques := cc.DB().GetQuestions()
	if len(r) != len(ques) {
		return model.NewApiError(model.ApiErrorCount, "Question and answer counts not equal", nil)
	}

	score, successRate := cc.DB().CalculateScoreAndDegree(c.Param("username"), r)

	return service.ApiResponse(c, http.StatusAccepted, map[string]interface{}{"username": c.Param("username"), "score": score, "success_rate": successRate})
}
