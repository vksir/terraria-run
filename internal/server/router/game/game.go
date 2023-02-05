package game

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"terraria-run/internal/controller"
	"terraria-run/internal/server/model/common/commonresp"
	"terraria-run/internal/server/model/game/gameresp"
)

func LoadRouters(g *gin.RouterGroup) {
	g.GET("/game/players", getPlayers)
}

// players godoc
// @Summary			查看当前玩家
// @Tags			game
// @Accept			json
// @Produce			json
// @Success			200 {object} gameresp.Players
// @Failure			500 {object} commonresp.Err
// @Router			/game/players [get]
func getPlayers(c *gin.Context) {
	if controller.Status != controller.StatusActive {
		c.JSON(http.StatusBadRequest, commonresp.Err{
			Detail: fmt.Sprintf("Status is %s, cannot get players", controller.Status),
		})
		return
	}
	out, err := controller.Agent.RunCmd("playing")
	if err != nil {
		c.JSON(http.StatusInternalServerError, commonresp.Err{Detail: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gameresp.Players{Players: []string{out}})
}
