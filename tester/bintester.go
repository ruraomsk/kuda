package tester

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ruraomsk/kuda/hardware/bin"
)

func BinTest() {
	c, err := bin.LoadBin("test.bin")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// logger.Debug.Printf("%v", c)
	b, err := json.MarshalIndent(&c, "", "   ")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = ioutil.WriteFile("test.json", b, 0644)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = c.SaveBin("save.bin")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
