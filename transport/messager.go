package transport

import (
	"github.com/ruraomsk/ag-server/pudge"
	"github.com/ruraomsk/kuda/brams"
	"github.com/ruraomsk/kuda/data"
)

type Message struct {
	Messages map[string][]byte `json:"messages"`
}

// var emptyMessage = Message{}
var execute = map[string]interface{}{
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

func workMessage(message Message) (Message, bool) {

	return statusMessage(), true
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
