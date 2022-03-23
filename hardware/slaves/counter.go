package slaves

import (
	"fmt"
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/kuda/setup"
	"github.com/tbrandon/mbserver"
)

func StartCounter() {
	countserver := mbserver.NewServer()
	con := fmt.Sprintf(":%d", setup.Set.Counter.SPort)
	for i := 0; i < 3; i++ {
		countserver.HoldingRegisters[i] = 0
	}
	err := countserver.ListenTCP(con)
	if err != nil {
		logger.Info.Printf("Counter slave  %s", err.Error())
		return
	}

	logger.Info.Printf("Counter slave ready %s", con)
	for {
		time.Sleep(time.Second)
	}
}
