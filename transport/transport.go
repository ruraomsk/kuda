package transport

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"net"
	"time"

	"github.com/ruraomsk/TLServer/logger"
)

type Message struct {
	Messages map[string][]byte `json:"messages"`
}

var key []byte
var emptyMessage = Message{}

func ConnectWithServer(serverIP string) (net.Conn, error) {

	socket, err := net.Dial("tcp", serverIP)
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
	message, err := reader.ReadString('\n')
	if err != nil {
		logger.Error.Printf("Чтение ключа от %s %s", socket.RemoteAddr(), err.Error())
		socket.Close()
		return nil, err
	}
	key, err = base64.StdEncoding.DecodeString(message)
	if err != nil {
		logger.Error.Printf("Чтение ключа от %s %s", socket.RemoteAddr(), err.Error())
		socket.Close()
		return nil, err

	}
	return socket, nil
}

func GetMessageFromServer(socket net.Conn, inchan chan Message, toutin time.Duration) {
	defer socket.Close()
	defer close(inchan)
	reader := bufio.NewReader(socket)
	for {
		socket.SetReadDeadline(time.Now().Add(toutin))
		message, err := reader.ReadString('\n')
		if err != nil {
			logger.Error.Printf("Чтение сообщения от %s %s", socket.RemoteAddr().String(), err.Error())
			inchan <- emptyMessage
			return
		}
		mess, err := decode(message)
		if err != nil {
			logger.Error.Printf("Декодирование сообщения от %s %s", socket.RemoteAddr().String(), err.Error())
			inchan <- emptyMessage
			return
		}
		var inm Message
		err = json.Unmarshal(mess, &inm)
		if err != nil {
			logger.Error.Printf("Unmarshal  сообщения %v %s", message, err.Error())
			inchan <- emptyMessage
			return
		}
		inchan <- inm
	}
}

func SendMessageToServer(socket net.Conn, outchan chan Message, toutsend time.Duration) {
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
		socket.SetWriteDeadline(time.Now().Add(toutsend))
		str, err := code(buffer)
		if err != nil {
			logger.Error.Printf("Кодирование сообщения %v %s", buffer, err.Error())
			return
		}
		_, _ = writer.WriteString(str)
		_, _ = writer.WriteString("\n")
		err = writer.Flush()
		if err != nil {
			logger.Error.Printf("Передача сообщения для %s %s", socket.RemoteAddr().String(), err.Error())
			return
		}
	}
}
