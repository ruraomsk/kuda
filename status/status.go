package status

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/kuda/brams"
)

/*
	Подсистема хранения статуса устройства и плюс записи журнала
*/
const (
	Server = iota
	Hard
	Controller
	Technology
	NetWare
)

type Buffer struct {
	Name     string   `json:"name"`
	Statuses []Status `json:"stats"`
}

type Message struct {
	Type    int
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}
type Status struct {
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}

var (
	names = map[int]string{
		Server:     "Сервер",
		Hard:       "Оборудование",
		Controller: "Контроллер",
		Technology: "Технология",
		NetWare:    "Браузер",
	}
	lastStatuses map[string]Status
	Messages     chan Message
	CleanAll     chan interface{}
	Cap          = 100 //Емкость циклических буферов
	err          error
	db           *brams.Db
)

func appendMessage(name string, status Status) {
	var b Buffer
	buf, _ := db.ReadRecord(name)
	json.Unmarshal(buf, &b)
	if len(b.Statuses) < Cap {
		b.Statuses = append(b.Statuses, status)
	} else {
		b.Statuses = b.Statuses[1:]
		b.Statuses = append(b.Statuses, status)
	}
	db.WriteJSON(b)
}
func mainLoop() {
	// tick := time.NewTicker(10 * time.Second)
	for {
		select {
		case mess := <-Messages:
			name, ok := names[mess.Type]

			if !ok {
				logger.Error.Printf("Dismiss type %v", mess)
				continue
			}
			lastm, ok := lastStatuses[name]
			if !ok {
				s := new(Status)
				lastStatuses[name] = *s
			}
			if strings.Compare(lastm.Message, mess.Message) != 0 {
				s := new(Status)
				s.Time = mess.Time
				s.Message = mess.Message
				lastStatuses[name] = *s
				appendMessage(name, *s)
			}
		case <-CleanAll:
			lastStatuses = make(map[string]Status)
			brams.Drop("statuses")
			brams.CreateDb("statuses", "name")
			db, _ = brams.Open("statuses")
			for _, n := range names {
				buffer := new(Buffer)
				buffer.Name = n
				buffer.Statuses = make([]Status, 0, Cap)
				db.WriteJSON(buffer)
			}
			// case <-tick.C:
			// 	for _, n := range names {
			// 		buf, _ := db.ReadRecord(n)
			// 		var b Buffer
			// 		json.Unmarshal(buf, &b)
			// 		logger.Info.Printf("%v", b)
			// 	}

		}

	}
}

func StartStatus() error {
	lastStatuses = make(map[string]Status)
	Messages = make(chan Message, 100)
	db, err = brams.Open("statuses")
	if err != nil {
		//Первый запуск создаем базу данных и инициализируемся
		if err = brams.CreateDb("statuses", "name"); err != nil {
			return fmt.Errorf("Not building db statuses")
		}
		db, _ = brams.Open("statuses")
		for _, n := range names {
			buffer := new(Buffer)
			buffer.Name = n
			buffer.Statuses = make([]Status, 0, Cap)
			db.WriteJSON(buffer)
		}
	} else {
		for _, n := range names {
			buf, _ := db.ReadRecord(n)
			var v Buffer
			json.Unmarshal(buf, &v)
			if len(v.Statuses) > 0 {
				mes := v.Statuses[len(v.Statuses)-1]
				lastStatuses[n] = Status{Time: mes.Time, Message: mes.Message}
			}
		}
	}
	go mainLoop()
	return nil
}
