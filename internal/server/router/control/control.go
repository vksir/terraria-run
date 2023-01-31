package control

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"terraria-run/internal/server/model/action/actionreq"
	"terraria-run/internal/server/model/common/commonresp"
)

func LoadRouters(g *gin.RouterGroup) {
	g.POST("/control/:action", control)
}

// control godoc
// @Summary			服务器控制
// @Tags			control
// @Accept			json
// @Produce			json
// @Param			action path string true "[ start | stop | restart ]"
// @Success			200 {object} commonresp.Ok
// @Failure			500 {object} commonresp.Err
// @Router			/control/{action} [post]
func control(c *gin.Context) {
	params := actionreq.Params{}
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, commonresp.Err{Detail: err.Error()})
		return
	}
	switch params.Action {
	case "start":
	case "stop":
	case "restart":
	default:
		c.JSON(http.StatusBadRequest, commonresp.Err{Detail: "Invalid action"})
		return
	}
	c.JSON(http.StatusOK, commonresp.Ok{})
}
