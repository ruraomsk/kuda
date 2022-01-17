package transport

import (
	"bufio"
	"encoding/json"
	"net"
	"time"

	"github.com/ruraomsk/TLServer/logger"
)

type Message struct {
	Messages map[string][]byte `json:"messages"`
}

var secretCode []byte
var emptyMessage = Message{}

func ConnectWithServer(adress string) (net.Conn, error) {

	socket, err := net.Dial("tcp", adress)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(socket)
	writer := bufio.NewWriter(socket)
	socket.SetDeadline(time.Now().Add(time.Second * 10))
	writer.WriteString(messageConnect())
	writer.WriteString("\n")
	err = writer.Flush()
	if err != nil {
		logger.Error.Printf("Передача сообщения для %s %s", socket.RemoteAddr(), err.Error())
		socket.Close()
		return nil, err
	}
	message, err := reader.ReadBytes('\n')
	if err != nil {
		logger.Error.Printf("Чтение сообщения от %s %s", socket.RemoteAddr(), err.Error())
		socket.Close()
		return nil, err
	}
	secretCode = message
	return socket, nil
}

func GetMessageFromServer(socket net.Conn, inchan chan Message, toutin *time.Duration) {
	defer socket.Close()
	defer close(inchan)
	reader := bufio.NewReader(socket)
	for {
		socket.SetReadDeadline(time.Now().Add(*toutin))
		message, err := reader.ReadBytes('\n')
		if err != nil {
			logger.Error.Printf("Чтение сообщения от %s %s", socket.RemoteAddr(), err.Error())
			inchan <- emptyMessage
			return
		}
		message, err = decode(message)
		if err != nil {
			logger.Error.Printf("Декодирование сообщения от %s %s", socket.RemoteAddr(), err.Error())
			inchan <- emptyMessage
			return
		}
		var inm Message
		err = json.Unmarshal(message, &inm)
		if err != nil {
			logger.Error.Printf("Unmarshal  сообщения %v %s", message, err.Error())
			inchan <- emptyMessage
			return
		}
		inchan <- inm
	}
}

func SendMessageToServer(socket net.Conn, outchan chan Message, toutsend *time.Duration) {
	defer socket.Close()
	defer close(outchan)
	writer := bufio.NewWriter(socket)
	for {
		message := <-outchan
		buffer, err := json.Marshal(message)
		if err != nil {
			logger.Error.Printf("Marshal  сообщения %v %s", message, err.Error())
			return
		}
		socket.SetWriteDeadline(time.Now().Add(*toutsend))
		buffer, err = code(buffer)
		if err != nil {
			logger.Error.Printf("Кодирование сообщения %v %s", buffer, err.Error())
			return
		}
		_, _ = writer.WriteString(string(buffer))
		_, _ = writer.WriteString("\n")
		err = writer.Flush()
		if err != nil {
			logger.Error.Printf("Передача сообщения для %s %s", socket.RemoteAddr(), err.Error())
			return
		}
	}
}