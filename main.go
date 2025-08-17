package main

import (
	"go-srv-bootstrap/config"
	"go-srv-bootstrap/helpers"
	"go-srv-bootstrap/server"
)

var masterCtx *helpers.ProcessContext

func init() {
	helpers.InitProcessContext()
}

func main() {
	masterCtx = helpers.GetProcessContext()
	masterCtx.AddWorker(1)
	srv := server.GetServer()

	go server.RunServer(masterCtx)
	config.GetLogger().Println("Server is running at ", srv.Addr)

	masterCtx.WaitForFinish()
}
