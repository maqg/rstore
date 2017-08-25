package api

import (
	"octlink/rstore/utils"

	"github.com/gin-gonic/gin"
)

// Router for api management
func (api *API) Router() *gin.Engine {

	var baseDir = ""

	router := gin.New()

	gin.SetMode(gin.ReleaseMode)

	exist := utils.IsFileExist("frontend")
	if !exist {
		baseDir = "../"
	}

	router.LoadHTMLGlob(baseDir + "frontend/apitest/templates/*.html")
	router.Static("/static", baseDir+"frontend/static")

	router.GET("/api/test/", api.LoadTestPage)

	router.GET("/api/", api.Test)
	router.POST("/api/", api.Dispatch)

	return router
}
