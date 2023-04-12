package controllers

import (
	"botgpt/internal/enum"
	"botgpt/internal/utils"
	"botgpt/internal/utils/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type StatusController struct {
}

func NewStatusController() *StatusController {
	return &StatusController{}
}

func (h StatusController) Status(c *gin.Context) {
	statusInfo := utils.StatusInfo

	result := response.Make(true, enum.SUCCESS, "success", statusInfo)

	utils.SendResponse(http.StatusOK, result, c)
}
