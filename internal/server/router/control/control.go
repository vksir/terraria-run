package control

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"terraria-run/internal/controller"
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
		if status := controller.Status(); status != controller.StatusInactive {
			c.JSON(http.StatusBadRequest, commonresp.Err{
				Detail: fmt.Sprintf("Status is %s, cannot start", status),
			})
			return
		}
		if err := controller.Start(); err != nil {
			c.JSON(http.StatusForbidden, commonresp.Err{Detail: err.Error()})
			return
		}
	case "stop":
		if status := controller.Status(); status != controller.StatusActive && status != controller.StatusStarting {
			c.JSON(http.StatusBadRequest, commonresp.Err{
				Detail: fmt.Sprintf("Status is %s, cannot stop", status),
			})
			return
		}
		if err := controller.Stop(); err != nil {
			c.JSON(http.StatusForbidden, commonresp.Err{Detail: err.Error()})
			return
		}
	case "restart":
		// TODO
	default:
		c.JSON(http.StatusBadRequest, commonresp.Err{Detail: fmt.Sprintf("Invalid action: %s", params.Action)})
		return
	}
	c.JSON(http.StatusOK, commonresp.Ok{})
}
