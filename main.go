package main

import "github.com/dhinesh20022806/node-s3/app/cmd/server"

func main() {

	server.GrpcRun()
	server.Run()

}
