package main

import (
	"github.com/jenujari/go-srv-bootstrap/config"
	"github.com/jenujari/go-srv-bootstrap/helpers"
	"github.com/jenujari/go-srv-bootstrap/server"
)

var cmder *helpers.Commander

func init() {
	cmder = helpers.NewCommander()
}

func main() {
	cmder.AddWorker(1)
	srv := server.GetServer()

	go server.RunServer(cmder)
	config.GetLogger().Println("Server is running at ", srv.Addr)

	cmder.WaitForFinish()
}
