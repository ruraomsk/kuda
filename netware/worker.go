package netware

import (
	"bufio"
	"net"

	"github.com/ruraomsk/TLServer/logger"
	"github.com/ruraomsk/kuda/commander"
)

func workerNetware(socket net.Conn) {
	defer socket.Close()
	reader := bufio.NewReader(socket)
	writer := bufio.NewWriter(socket)
	for {
		in, err := reader.ReadString('\n')
		if err != nil {
			logger.Error.Printf("Чтение команды для %s %s", socket.RemoteAddr(), err.Error())
			return
		}
		logger.Info.Printf("%s", in)
		result, err := commander.DoCommand(in)
		if err != nil {
			result = err.Error()
		}
		writer.WriteString(result)
		writer.WriteString("\n")
		writer.Flush()
	}
}
