package server

import (
	"github.com/dhinesh20022806/node-s3/app/internal/services"
	"github.com/gin-gonic/gin"
)



func Run() {


	server := gin.Default()

	server.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "success",
		})
	})

	


	server.Run(":8999")

}

func GrpcRun() {
	serviceGrpc := services.NewGrpcServer(":50052")

	serviceGrpc.Run()
}