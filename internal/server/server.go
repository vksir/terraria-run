package server

import (
	"github.com/gin-gonic/gin"
	"log"
	"terraria-run/internal/server/router/control"
	"terraria-run/internal/server/router/mod"
	"terraria-run/internal/server/swagger"
)

// @title           Terraria Run
// @version         1.0

func Run() {
	e := gin.Default()
	loadRouters(e)
	err := e.Run("0.0.0.0:5778")
	if err != nil {
		log.Panic(err)
	}
}

func loadRouters(e *gin.Engine) {
	publicGroup := e.Group("")
	control.LoadRouters(publicGroup)
	mod.LoadRouters(publicGroup)
	swagger.LoadRouters(publicGroup)
}
