package brams

import (
	"io/ioutil"
	"strings"
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/kuda/setup"
)

func saveDBs() {
	dbs.Lock()
	defer dbs.Unlock()
	for _, db := range dbs.dbs {
		err := db.saveToFile()
		if err != nil {
			logger.Error.Printf("Saving БД %s %s", db.Name, err.Error())
		}
	}
}

func workerDBFs(stop chan interface{}) {
	ticker := time.NewTicker(time.Duration(sbrams.Step) * time.Second)
	for {
		select {
		case <-ticker.C:
			saveDBs()
		case <-stop:
			work = false
			saveDBs()
			return
		}
	}
}

func getListFilesDbs() []string {
	list := make([]string, 0)
	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		return list
	}
	for _, dir := range dirs {
		if dir.IsDir() {
			continue
		}
		if strings.HasSuffix(dir.Name(), ext) {

			list = append(list, strings.TrimSuffix(dir.Name(), ext))
		}
	}
	return list
}
func StartBrams(dbstop chan interface{}) error {
	path = setup.Set.SetupBrams.DbPath

	work = true
	for _, db := range getListFilesDbs() {
		if err := addDbFromJson(db); err != nil {
			return err
		}
	}
	go workerDBFs(dbstop)
	return nil
}
