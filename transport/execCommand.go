package transport

import (
	"encoding/json"

	"github.com/ruraomsk/ag-server/comm"
	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/ag-server/pudge"
	"github.com/ruraomsk/kuda/brams"
	"github.com/ruraomsk/kuda/data"
)

const (
	Control     = 2
	Sfdk        = 4
	SetPK       = 5
	SetCK       = 6
	SetNK       = 7
	SetDU       = 9
	SetDU2      = 10
	GetPrivDate = 12
	GetPriv     = 11
)

func execCommand(cmd comm.CommandARM) (Message, bool) {
	var err error
	var db, dbb, dbd *brams.Db
	logger.Info.Printf("команда %v", cmd)
	switch cmd.Command {
	case Control:
		if cmd.Params == 1 {
			logger.Info.Print("Начинаем прием привязок")
		} else {
			logger.Info.Print("Прием привязок закончен")
		}

	case Sfdk:
		db, err = brams.Open("StatusCommandDU")
		if err != nil {
			logger.Error.Printf("StatusCommandDU %s", err.Error())
			return emptyMessage, false
		}
		var scd pudge.StatusCommandDU
		db.ReadJSON(&scd)
		scd.IsReqSFDK1 = cmd.Params&1 != 0
		scd.IsReqSFDK2 = cmd.Params&2 != 0
		db.WriteJSON(scd)
		db.Close()

	case SetPK:
		db, err = brams.Open("StatusCommandDU")
		if err != nil {
			logger.Error.Printf("StatusCommandDU %s", err.Error())
			return emptyMessage, false
		}
		dbb, err = brams.Open("base")
		if err != nil {
			logger.Error.Printf("Base %s", err.Error())
			return emptyMessage, false
		}
		var scd pudge.StatusCommandDU
		var b data.BaseCtrl
		db.ReadJSON(&scd)
		dbb.ReadJSON(&b)
		b.PK = cmd.Params
		if cmd.Params == 0 {
			scd.IsPK = false
		} else {
			scd.IsPK = true
		}
		db.WriteJSON(scd)
		dbb.WriteJSON(b)
		db.Close()
		dbb.Close()
	case SetCK:
		db, err = brams.Open("StatusCommandDU")
		if err != nil {
			logger.Error.Printf("StatusCommandDU %s", err.Error())
			return emptyMessage, false
		}
		dbb, err = brams.Open("base")
		if err != nil {
			logger.Error.Printf("Base %s", err.Error())
			return emptyMessage, false
		}
		var scd pudge.StatusCommandDU
		var b data.BaseCtrl
		db.ReadJSON(&scd)
		dbb.ReadJSON(&b)
		b.CK = cmd.Params
		if cmd.Params == 0 {
			scd.IsCK = false
		} else {
			scd.IsCK = true
		}
		db.WriteJSON(scd)
		dbb.WriteJSON(b)
		db.Close()
		dbb.Close()
	case SetNK:
		db, err = brams.Open("StatusCommandDU")
		if err != nil {
			logger.Error.Printf("StatusCommandDU %s", err.Error())
			return emptyMessage, false
		}
		dbb, err = brams.Open("base")
		if err != nil {
			logger.Error.Printf("Base %s", err.Error())
			return emptyMessage, false
		}
		var scd pudge.StatusCommandDU
		var b data.BaseCtrl
		db.ReadJSON(&scd)
		dbb.ReadJSON(&b)
		b.NK = cmd.Params
		if cmd.Params == 0 {
			scd.IsNK = false
		} else {
			scd.IsNK = true
		}
		db.WriteJSON(scd)
		dbb.WriteJSON(b)
		db.Close()
		dbb.Close()
	case SetDU:
		db, err = brams.Open("StatusCommandDU")
		if err != nil {
			logger.Error.Printf("StatusCommandDU %s", err.Error())
			return emptyMessage, false
		}
		dbd, err = brams.Open("DK")
		if err != nil {
			logger.Error.Printf("DK %s", err.Error())
			return emptyMessage, false
		}
		var scd pudge.StatusCommandDU
		var dk pudge.DK
		db.ReadJSON(&scd)
		dbd.ReadJSON(&dk)
		if cmd.Params == 0 {
			dk.RDK = 8
			scd.IsDUDK1 = false
		} else {
			scd.IsDUDK1 = true
			dk.FDK = cmd.Params
			dk.RDK = 4
			dk.DDK = 6
		}
		db.WriteJSON(scd)
		dbd.WriteJSON(dk)
		db.Close()
		dbd.Close()
	case GetPrivDate:
		db, err = brams.Open("dates")
		if err != nil {
			err := brams.CreateDb("dates", "name")
			if err != nil {
				logger.Error.Printf("Не могу создать dates %s", err.Error())
				return emptyMessage, false
			}
			db, _ = brams.Open("dates")
			for _, d := range data.DatesList {
				db.WriteJSON(d)
			}
		}
		var dt data.Dates
		m := emptyMessage
		mas := make([]data.Dates, 0)
		dts, _ := db.ReadListKeys(0)
		for _, dn := range dts {
			db.ReadJSON(&dt, dn)
			mas = append(mas, dt)
		}
		buf, _ := json.Marshal(mas)
		m.Messages["dates"] = buf
		db.Close()
		return m, true
	case GetPriv:

	}
	return emptyMessage, false
}
