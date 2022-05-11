package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary HeartBeat
// @produce text/plain
// @Success 200 string string "Alive"
// @Router /heartBeat [get]
func HeartBeat(c *gin.Context) {
	c.String(http.StatusOK, "Alive")
}
