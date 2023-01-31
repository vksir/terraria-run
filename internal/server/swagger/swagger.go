package swagger

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	_ "terraria-run/docs"
)

func LoadRouters(g *gin.RouterGroup) {
	g.GET("/docs", redirect)
	g.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func redirect(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/docs/index.html")
}
