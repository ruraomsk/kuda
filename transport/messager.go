package transport

import (
	"encoding/json"

	"github.com/ruraomsk/ag-server/comm"
	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/ag-server/pudge"
	"github.com/ruraomsk/kuda/brams"
	"github.com/ruraomsk/kuda/data"
)

type Message struct {
	Messages map[string][]byte `json:"messages"`
}

var (
	emptyMessage = Message{}
	execute      = map[string]interface{}{
		"base":            data.BaseCtrl{},
		"traffic":         pudge.Traffic{},
		"Status":          pudge.Status{},
		"StatusCommandDU": pudge.StatusCommandDU{},
		"DK":              pudge.DK{},
		"Model":           pudge.Model{},
		"ErrorDevice":     pudge.ErrorDevice{},
		"GPS":             pudge.GPS{},
		"Input":           pudge.Input{},
	}
)

func workMessage(message Message) (Message, bool) {
	if _, is := message.Messages["status"]; is {
		logger.Info.Print("Request status")
		return statusMessage(), true
	}
	if cmd, is := message.Messages["command"]; is {
		var command comm.CommandARM
		err := json.Unmarshal(cmd, &command)
		if err != nil {
			logger.Error.Printf("Encoding %v %s", string(cmd), err.Error())
			return emptyMessage, false
		}
		if message, is := execCommand(command); is {
			return message, is
		}
		return emptyMessage, false
	}
	return emptyMessage, false
}
func statusMessage() Message {
	m := Message{Messages: make(map[string][]byte)}
	for name, buffer := range execute {
		db, err := brams.Open(name)
		if err != nil {
			brams.CreateDb(name)
			db, _ = brams.Open(name)
			db.WriteJSON(buffer)
		}
		buf, err := db.ReadRecord()
		if err != nil {
			continue
		}
		m.Messages[name] = buf
	}
	return m
}
func exitMessage() Message {
	m := Message{Messages: make(map[string][]byte)}
	db, err := brams.Open("base")
	if err != nil {
		return emptyMessage
	}
	buf, _ := db.ReadRecord()
	db.Close()
	m.Messages["ErrorDevice"] = buf
	return m
}
