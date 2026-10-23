package main

import (
	"github.com/dhinesh20022806/node-s3/app/cmd/server"
	_ "github.com/dhinesh20022806/node-s3/app/pkg/dispatcher"
)

func main() {

	server.GrpcRun()
	server.Run()

}
