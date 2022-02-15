package netware

import (
	"fmt"
	"net"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/kuda/setup"
	"github.com/ruraomsk/kuda/status"
)

func listenCommand() {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", setup.Set.Netware.Port))
	if err != nil {
		logger.Error.Printf("Error opening port %d", setup.Set.Netware.Port)
		return
	}
	for {
		socket, err := ln.Accept()
		if err != nil {
			logger.Error.Printf("Accept error %s", err.Error())

		}
		go workerNetware(socket)
	}
}

func StartNetware() {
	logger.Info.Println("Netware start")
	status.NetwareMessage("Internal Http server start")
	go listenCommand()
	go server()
}
