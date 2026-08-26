package main

import (
	_ "gotask/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func registerSwaggerRoutes(router *gin.Engine, username, password string) {
	swagger := router.Group("/swagger", gin.BasicAuth(gin.Accounts{
		username: password,
	}))
	swagger.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
