package netware

import (
	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/kuda/brams"
)

/*
	Программа для связи с сервером верхнего уровня

*/

func server() {
	db, err := brams.Open("setup")
	if err != nil {
		logger.Info.Println("setup netware not found")
		return
	}
	db.Close()
}
