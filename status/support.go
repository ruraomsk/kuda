package status

import "time"

func ServerMessage(message string) {
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
