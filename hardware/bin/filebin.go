package bin

import (
	"io/ioutil"

	"github.com/ruraomsk/ag-server/logger"
)

var (
	buff []byte
	err  error
	nm   int
)

func LoadBin(path string) (*CMK, error) {
	c := new(CMK)
	buff, err = ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c.BaseOption = int(buff[0])
	c.BaseRPU = int(buff[1])
	c.NumTU = int(buff[2])
	c.NumDK = int(buff[3])
	pos := 4
	for nm = 1; nm < 30; nm++ {
		ptr := int(buff[pos]) << 8
		ptr |= int(buff[pos+1])
		ptr -= 0x8000
		c.convertToStruct(ptr, int(buff[pos+2]), int(buff[pos+3]))
		pos += 4
	}
	return c, nil
}
func (c *CMK) SaveBin(path string) error {
	buffer := make([]byte, 120)
	buffer[0] = byte(c.BaseOption)
	buffer[1] = byte(c.BaseRPU)
	buffer[2] = byte(c.NumTU)
	buffer[3] = byte(c.NumDK)
	ptr := 4
	for i := 0; i < 29; i++ {
		b, _, count := c.toBin(i + 1)
		adr := len(buffer) + 0x8000
		buffer[ptr] = byte(adr >> 8)
		buffer[ptr+1] = byte(adr & 0xff)
		buffer[ptr+2] = byte(len(b) / count)
		buffer[ptr+3] = byte(count)
		logger.Debug.Printf("newMass %d %x %d %d", i+1, adr, int(len(b)/count), count)
		buffer = append(buffer, b...)
		ptr += 4
	}
	err := ioutil.WriteFile(path, buffer, 0644)
	return err
}
