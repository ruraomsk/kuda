package brams

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"os"
	"sync"

	"github.com/ruraomsk/kuda/setup"
)

var (
	dbs struct {
		sync.RWMutex
		dbs map[string]*Db
	}
	path               string
	ErrKeyNotFound     = errors.New("error: key not found")
	ErrKeyBad          = errors.New("error: key bad")
	ErrWrongParameters = errors.New("error: wrong param")
	ErrStopped         = errors.New("error: brams stopped")
	ext                = ".json"
	extData            = ".dat"
	sbrams             setup.SetupBrams
	work               bool
	oneKey             = "oneKey"
)

const (
	JSON = "json"
)

//Db Структура для доступа к БД
type Db struct {
	sync.RWMutex
	Name   string           `json:"name"`
	Defkey []string         `json:"dk"`
	values map[string]Value `json:"-"`
	UID    uint64           `json:"uid"`
	update bool
	fs     bool
}
type Value struct {
	UID   uint64 `json:"uid"`
	Value []byte `json:"v"`
}

func init() {
	dbs.dbs = make(map[string]*Db)
	path = "./"
	sbrams.Step = 1
	work = false
}

func (db *Db) Close() {
	db.RWMutex.Lock()
	defer db.RWMutex.Unlock()
}
func (db *Db) makeFullKey(keys []interface{}) (string, error) {
	var err error
	if len(keys) == 0 {
		return oneKey, nil
	}
	full := new(bytes.Buffer)
	for _, v := range keys {
		switch vt := v.(type) {
		case []byte:
			_, err = full.Write(vt)
		case string:
			_, err = full.WriteString(vt)
		default:
			err = binary.Write(full, binary.BigEndian, v)

		}
		if err != nil {
			return "", err
		}
	}
	return full.String(), nil
}
func (db *Db) makeFullKeyOnValue(value []byte) (string, error) {
	if len(db.Defkey) == 0 {
		return oneKey, nil
	}
	m := make(map[string]interface{})
	err := json.Unmarshal(value, &m)
	if err != nil {
		return "", err
	}
	full := new(bytes.Buffer)
	for _, kn := range db.Defkey {
		v, is := m[kn]
		if !is {
			return "", ErrKeyBad
		}
		switch vt := v.(type) {
		case []byte:
			_, err = full.Write(vt)
		case string:
			_, err = full.WriteString(vt)
		default:
			err = binary.Write(full, binary.BigEndian, v)

		}
		if err != nil {
			return "", err
		}
	}
	return full.String(), nil
}

func (db *Db) saveToFile() error {
	db.RWMutex.Lock()
	defer db.RWMutex.Unlock()
	if !db.fs {
		db.update = false
		return nil
	}
	if !db.update {
		return nil
	}
	// logger.Info.Printf("Сохраняем БД %s", db.Name)
	fname := path + db.Name + ext
	buffer, err := json.Marshal(db)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fname, buffer, os.FileMode(0644))
	if err != nil {
		return err
	}
	fname = path + db.Name + extData
	temp := make([]Value, 0)
	for _, t := range db.values {
		temp = append(temp, t)
	}
	buffer, err = json.Marshal(temp)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fname, buffer, os.FileMode(0644))
	if err != nil {
		return err
	}

	db.update = false
	return nil
}
func (db *Db) makeListKeys(limit int, keys []interface{}) ([]string, error) {
	if limit == 0 {
		limit = math.MaxInt32
	}
	if len(db.Defkey) == 0 {
		return make([]string, 0), nil
	}
	var err error
	m := make(map[string]interface{})
	result := make([]string, 0)
	if len(keys) == 0 {
		// logger.Debug.Printf("Start len keys=0 %v", db.values)
		for full := range db.values {
			if limit > 0 {
				// logger.Debug.Printf("append key %s", full)
				result = append(result, full)
				limit--
			}

		}
		return result, nil
	}
	key := make([][]byte, 0)
	count := 0
	for _, k := range keys {
		nv := new(bytes.Buffer)
		switch kt := k.(type) {
		case nil:
			key = append(key, make([]byte, 0))
			continue
		case []byte:
			_, err = nv.Write(kt)
			count++
		case string:
			_, err = nv.WriteString(kt)
			count++
		case int:
			val := float64(kt)
			err = binary.Write(nv, binary.BigEndian, val)
			count++
		default:
			err = binary.Write(nv, binary.BigEndian, k)
			count++
		}
		key = append(key, nv.Bytes())
		if err != nil {
			return make([]string, 0), err
		}
	}
	for _, val := range db.values {
		err = json.Unmarshal(val.Value, &m)
		if err != nil {
			return make([]string, 0), err
		}
		found := 0
		for i, k := range db.Defkey {
			if len(key[i]) == 0 {
				continue
			}
			v, is := m[k]
			if !is {
				return make([]string, 0), ErrKeyBad
			}
			vv := new(bytes.Buffer)
			switch vt := v.(type) {
			case []byte:
				_, err = vv.Write(vt)
			case string:
				_, err = vv.WriteString(vt)
			default:
				err = binary.Write(vv, binary.BigEndian, v)
			}
			if err != nil {
				return make([]string, 0), err
			}
			if bytes.Equal(key[i], vv.Bytes()) {
				found++
			}
		}
		if (found == count) && limit > 0 {
			full, err := db.makeFullKeyOnValue(val.Value)
			if err != nil {
				return make([]string, 0), err
			}
			result = append(result, full)
			limit--
		}

	}
	return result, nil
}
