package transport

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/ag-server/pudge"
	"github.com/ruraomsk/kuda/brams"
	"github.com/ruraomsk/kuda/status"
)

var key []byte

func ConnectWithServer(serverIP string, id int) (net.Conn, error) {

	socket, err := net.Dial("tcp", serverIP)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(socket)
	writer := bufio.NewWriter(socket)
	socket.SetDeadline(time.Now().Add(time.Second * 10))
	writer.WriteString(fmt.Sprintf("%d", id))
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
	message = strings.ReplaceAll(message, "\n", "")
	if strings.Contains(message, "notFound!") {
		logger.Error.Printf("Не прописан на сервере")
		status.ServerMessage("Не прописан на сервере", 3)
		socket.Close()
		return nil, fmt.Errorf("не прописан на сервере")
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
			return
		}

		addTraffic(0, len(message))
		message = strings.ReplaceAll(message, "\n", "")
		mess, err := decode(message)
		if err != nil {
			logger.Error.Printf("Декодирование сообщения от %s %s", socket.RemoteAddr().String(), err.Error())
			return
		}
		var inm Message
		err = json.Unmarshal(mess, &inm)
		if err != nil {
			logger.Error.Printf("Unmarshal  сообщения %v %s", message, err.Error())
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
		n, _ := writer.WriteString(str)
		addTraffic(n, 0)
		_, _ = writer.WriteString("\n")
		err = writer.Flush()
		if err != nil {
			logger.Error.Printf("Передача сообщения для %s %s", socket.RemoteAddr().String(), err.Error())
			return
		}
	}
}
func addTraffic(in, out int) {
	var err error
	var db *brams.Db
	var tr = pudge.Traffic{}
	db, err = brams.Open("traffic")
	if err != nil {
		return
	}
	db.ReadJSON(&tr)
	tr.LastFromDevice1Hour += uint64(out)
	tr.LastToDevice1Hour += uint64(in)
	db.WriteJSON(tr)
	db.Close()
}
