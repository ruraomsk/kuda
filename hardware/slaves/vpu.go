package slaves

import (
	"fmt"
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/kuda/setup"
	"github.com/tbrandon/mbserver"
)

func StartVPU() {
	vpuserver := mbserver.NewServer()
	con := fmt.Sprintf(":%d", setup.Set.Vpu.SPort)
	for i := 0; i < 3; i++ {
		vpuserver.HoldingRegisters[i] = 0
	}
	err := vpuserver.ListenTCP(con)
	if err != nil {
		logger.Info.Printf("VPUr slave  %s", err.Error())
		return
	}
	logger.Info.Printf("VPUr slave ready %s", con)

	for {
		time.Sleep(time.Second)
	}
}
