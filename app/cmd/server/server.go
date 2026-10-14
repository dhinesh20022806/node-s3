package server

import "github.com/gin-gonic/gin"



func Run() {


	server := gin.Default()

	server.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "success",
		})
	})



	server.Run(":8999")

}