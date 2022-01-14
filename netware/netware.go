package netware

import (
	"fmt"
	"net"

	"github.com/ruraomsk/TLServer/logger"
	"github.com/ruraomsk/kuda/setup"
)

func listenCommand() {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", setup.Set.Netware.Port))
	if err != nil {
		logger.Error.Printf("Ошибка открытия порта %d", setup.Set.Netware.Port)
		return
	}
	for {
		socket, err := ln.Accept()
		if err != nil {
			logger.Error.Printf("Ошибка accept %s", err.Error())

		}
		go workerNetware(socket)
	}
}

func StartNetware() {
	logger.Info.Println("Netware start")
	go listenCommand()
}