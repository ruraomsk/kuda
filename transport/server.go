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
var ExitDevice chan interface{}
var ticker *time.Ticker
var IsConnected bool

func StartServerExchange(ip string) {
	ExitDevice = make(chan interface{})
	breakWork := false
	for {
		var base = data.BaseCtrl{ID: 10001, TimeOut: 500, TMax: 120}
		db, err := brams.Open("base")
		if err != nil {
			logger.Error.Print("db base not found... create")
			brams.CreateDb("base")
			db, _ = brams.Open("base")
			db.WriteJSON(base)
		}
		err = db.ReadJSON(&base)
		if err != nil {
			logger.Error.Printf("db base not reading %s", err.Error())
			status.ServerMessage(err.Error(), 11)
			time.Sleep(5 * time.Second)
			continue
		}
		// base.ID = 10001
		base.TimeDevice = time.Now()
		base.Base = true
		base.TMax = 120
		db.WriteJSON(base)
		db.Close()
		socket, err := ConnectWithServer(ip, base.ID)
		if err != nil {
			status.ServerMessage(err.Error(), 10)
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
		ticker = time.NewTicker(time.Duration(base.TMax) * time.Second)
		hour := time.NewTicker(time.Hour)
		oneSecond := time.NewTicker(time.Second)
		status.ServerMessage("Установлена связь с сервером", 0)
		IsConnected = true
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
					ticker.Reset(time.Duration(base.TMax) * time.Second)
				}

			case <-ticker.C:
				writeChan <- statusMessage()
				logger.Info.Printf("шлем статус %d", base.TMax)
				ticker.Reset(time.Duration(base.TMax) * time.Second)
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
			case <-ExitDevice:
				writeChan <- exitMessage()
				breakWork = true
				break loop
			case <-oneSecond.C:
				db, err := brams.Open("base")
				if err == nil {
					var bs data.BaseCtrl
					db.ReadJSON(&bs)
					bs.TimeDevice = time.Now()
					db.WriteJSON(bs)
					db.Close()
				}
			}

		}
		if breakWork {
			break
		}
		IsConnected = false
		socket.Close()
		status.ServerMessage("Отсутствует связь с сервером", 1)
		ticker.Stop()
		hour.Stop()
		time.Sleep(time.Second)
	}
	logger.Info.Print("Обмен с сервером прекращен")
}
