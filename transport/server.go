package transport

import (
	"encoding/json"
	"time"

	"github.com/ruraomsk/kuda/brams"
	"github.com/ruraomsk/kuda/data"
	"github.com/ruraomsk/kuda/status"
)

var ChangeChan chan string

func StartServerExchange(ip string) {

	for {
		socket, err := ConnectWithServer(ip)
		if err != nil {
			status.ServerMessage(err.Error())
			time.Sleep(10 * time.Second)
			continue
		}
		db, err := brams.Open("base")
		if err != nil {
			time.Sleep(10 * time.Second)
			continue
		}
		buf, err := db.ReadOneRecord()
		if err != nil {
			status.ServerMessage(err.Error())
			time.Sleep(10 * time.Second)
			continue
		}
		var base data.BaseCtrl
		_ = json.Unmarshal(buf, &base)
		readChan := make(chan Message)
		writeChan := make(chan Message)
		readTout := time.Duration(base.TimeOut * int64(time.Second))
		writeTout := time.Duration(10 * time.Second)
		go GetMessageFromServer(socket, readChan, readTout)
		go SendMessageToServer(socket, writeChan, writeTout)
		ticker := time.NewTicker(time.Duration(base.TMax * int64(time.Minute)))
	loop:
		for {
			select {
			case message, ok := <-readChan:
				if !ok {
					break loop
				}

			case <-ticker.C:

			case command, ok := <-ChangeChan:
				if !ok {
					break loop
				}

			}

		}

		close(writeChan)

	}
}
