package transport

import (
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/ag-server/pudge"
	"github.com/ruraomsk/kuda/brams"
	"github.com/ruraomsk/kuda/data"
	"github.com/ruraomsk/kuda/status"
)

var ChangeChan chan string
var ticker *time.Ticker

func StartServerExchange(ip string) {

	for {
		var base data.BaseCtrl
		db, err := brams.Open("base")
		if err != nil {
			logger.Error.Print("db base not found... create")
			brams.CreateDb("base")
			var base = data.BaseCtrl{ID: 10001, TimeOut: 500, TMax: 1}
			db, _ = brams.Open("base")
			db.WriteJSON(base)
		}
		err = db.ReadJSON(&base)
		if err != nil {
			status.ServerMessage(err.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		socket, err := ConnectWithServer(ip, base.ID)
		if err != nil {
			status.ServerMessage(err.Error())
			logger.Error.Printf("connect %s %s", ip, err.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		readChan := make(chan Message)
		writeChan := make(chan Message)
		readTout := time.Duration(base.TimeOut * int64(time.Second))
		writeTout := time.Duration(10 * time.Second)
		go GetMessageFromServer(socket, readChan, readTout)
		go SendMessageToServer(socket, writeChan, writeTout)
		ticker = time.NewTicker(time.Duration(base.TMax * int64(time.Minute)))
		hour := time.NewTicker(time.Hour)
		status.ServerMessage("Установлена связь с сервером")
	loop:
		for {
			select {
			case message, ok := <-readChan:
				if !ok {
					break loop
				}
				logger.Debug.Printf("Пришло %v", message)
				ReplayMessage, send := workMessage(message)
				if send {
					writeChan <- ReplayMessage
					ticker.Reset(time.Duration(base.TMax * int64(time.Minute)))
				}

			case <-ticker.C:
				db, _ := brams.Open("traffic")
				var tr pudge.Traffic
				db.ReadJSON(tr)
				logger.Info.Printf("traffic %v", tr)
				writeChan <- statusMessage()
				ticker.Reset(time.Duration(base.TMax * int64(time.Minute)))
			case <-hour.C:
				db, _ := brams.Open("traffic")
				var tr pudge.Traffic
				db.ReadJSON(&tr)
				tr.FromDevice1Hour = tr.LastFromDevice1Hour
				tr.LastFromDevice1Hour = 0
				tr.ToDevice1Hour = tr.LastToDevice1Hour
				tr.LastToDevice1Hour = 0
				db.WriteJSON(tr)
				db.Close()
			case command, ok := <-ChangeChan:
				if !ok {
					break loop
				}
				logger.Debug.Printf("Пришло %v", command)

			}

		}
		socket.Close()
		status.ServerMessage("Отсутсвует связь с сервером")
		ticker.Stop()
	}
}
