package status

import (
	"time"

	"github.com/ruraomsk/ag-server/pudge"
	"github.com/ruraomsk/kuda/brams"
)

func ServerMessage(message string, code int) {
	db, err := brams.Open("Status")
	if err == nil {
		var s pudge.Status
		db.ReadJSON(&s)
		s.StatusServer = code
		db.WriteJSON(s)
		db.Close()
	}
	Messages <- Message{Type: Server, Time: time.Now(), Message: message}
}
func HardMessage(message string) {
	Messages <- Message{Type: Hard, Time: time.Now(), Message: message}
}
func ControllerMessage(message string) {
	Messages <- Message{Type: Controller, Time: time.Now(), Message: message}
}
func TechMessage(message string) {
	Messages <- Message{Type: Technology, Time: time.Now(), Message: message}
}
func NetwareMessage(message string) {
	Messages <- Message{Type: NetWare, Time: time.Now(), Message: message}
}
