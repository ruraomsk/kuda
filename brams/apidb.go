package brams

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//Open открытие базы данных
//Открывает базу по имени и возвращает указатель на нее
//Если базы данных нет то возвращает ошибку
func Open(name string) (*Db, error) {
	if !work {
		return nil, ErrStopped
	}
	dbs.RLock()
	defer dbs.RUnlock()
	if db, ok := dbs.dbs[name]; ok {
		return db, nil
	}
	return nil, fmt.Errorf("need create db %s", name)
}
func (db *Db) Count() uint64 {
	db.RWMutex.RLock()
	defer db.RWMutex.RUnlock()
	return uint64(len(db.values))
}
func GetListDBNames() []string {
	if !work {
		return make([]string, 0)
	}
	dbs.RLock()
	defer dbs.RUnlock()
	res := make([]string, 0)
	for _, db := range dbs.dbs {
		res = append(res, db.Name)
	}
	return res
}

//Drop удаление базы данных
func Drop(name string) {
	if !work {
		return
	}
	dbs.Lock()
	defer dbs.Unlock()
	db, ok := dbs.dbs[name]
	if !ok {
		return
	}
	delete(dbs.dbs, name)
	if db.fs {
		fname := path + name + ext
		_ = os.Remove(fname)
		fname = path + name + extData
		_ = os.Remove(fname)
	}
}

//addDbFromJson добавляет бд в пул бд
func addDbFromJson(name string) error {
	dbs.Lock()
	defer dbs.Unlock()
	if _, ok := dbs.dbs[name]; ok {
		return fmt.Errorf("db %s is exist ", name)
	}
	fname := path + name + ext
	_, err := os.Stat(fname)
	if err != nil {
		return err
	}
	db := new(Db)
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, &db)
	if err != nil {
		return err
	}
	fname = path + name + extData
	_, err = os.Stat(fname)
	if err != nil {
		return err
	}
	var data []Value
	buf, err = ioutil.ReadFile(fname)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, &data)
	if err != nil {
		return err
	}
	var uid uint64
	db.values = make(map[string]Value)
	uid = 0
	for _, dt := range data {
		if uid < dt.UID {
			uid = dt.UID
		}
		fullkey, err := db.makeFullKeyOnValue(dt.Value)
		if err != nil {
			return err
		}
		db.values[fullkey] = dt
	}

	db.fs = true
	dbs.dbs[name] = db
	return nil
}

//CreateDb cоздает бд и присваивает описание ключа
// где defkey массив имен переменных из value json
func CreateDb(name string, defkeys ...string) error {
	if !work {
		return ErrStopped
	}
	dbs.Lock()
	defer dbs.Unlock()
	if _, ok := dbs.dbs[name]; ok {
		return fmt.Errorf("db %s is exist ", name)
	}
	if len(defkeys) == 0 {
		return ErrWrongParameters
	}
	db := new(Db)
	db.Name = name
	db.Defkey = defkeys
	db.values = make(map[string]Value)
	db.UID = 0
	db.fs = true
	db.update = true
	dbs.dbs[name] = db
	db.saveToFile()
	// fname := path + name + ext
	// os.Remove(fname)
	// // _, err := os.Stat(fname)
	// // if err == nil {
	// // 	return fmt.Errorf("db file %s is exist the path %s", name, path)
	// // }
	// buf, err := json.Marshal(&db)
	// if err != nil {
	// 	return err
	// }
	// err = ioutil.WriteFile(fname, buf, os.FileMode(0644))
	// if err != nil {
	// 	return err
	// }
	// os.Remove(path + name + extData)
	return nil
}
func CreateDbInMemory(name string, defkeys ...string) error {
	if !work {
		return ErrStopped
	}
	dbs.Lock()
	defer dbs.Unlock()
	if _, ok := dbs.dbs[name]; ok {
		return fmt.Errorf("db %s is exist ", name)
	}
	db := new(Db)
	db.Name = name
	db.Defkey = make([]string, 0)
	if len(defkeys) != 0 {
		db.Defkey = append(db.Defkey, defkeys...)
	}
	db.values = make(map[string]Value)
	db.fs = false
	db.UID = 0
	dbs.dbs[name] = db
	return nil
}
func (db *Db) WriteJSON(value interface{}) error {
	buf, _ := json.Marshal(value)
	return db.WriteRecord(buf)
}

func (db *Db) WriteRecord(value []byte) error {
	if !work {
		return ErrStopped
	}
	v := new(Value)
	v.Value = value
	fullkey, err := db.makeFullKeyOnValue(value)
	if err != nil {
		return err
	}
	db.RWMutex.RLock()
	old, is := db.values[fullkey]
	if !is {
		//Insert
		db.UID++
		v.UID = db.UID
	} else {
		//Replace
		v.UID = old.UID
	}
	db.RWMutex.RUnlock()
	db.RWMutex.Lock()
	db.values[fullkey] = *v
	db.update = true
	db.RWMutex.Unlock()
	return nil
}
func (db *Db) DeleteRecord(keys ...interface{}) error {
	if !work {
		return ErrStopped
	}
	if len(keys) != len(db.Defkey) {
		return ErrWrongParameters
	}
	full, err := db.makeFullKey(keys)
	if err != nil {
		return err
	}
	db.RWMutex.RLock()
	_, is := db.values[full]
	db.RWMutex.RUnlock()
	if !is {
		return ErrKeyNotFound
	}
	db.RWMutex.Lock()
	delete(db.values, full)
	db.RWMutex.Unlock()
	return nil
}

func (db *Db) ReadRecord(keys ...interface{}) ([]byte, error) {
	if !work {
		return make([]byte, 0), ErrStopped
	}
	if len(keys) != len(db.Defkey) {
		return make([]byte, 0), ErrWrongParameters
	}
	full, err := db.makeFullKey(keys)
	if err != nil {
		return make([]byte, 0), err
	}
	db.RWMutex.RLock()
	defer db.RWMutex.RUnlock()
	value, is := db.values[full]
	if !is {
		return make([]byte, 0), ErrKeyNotFound
	}
	return value.Value, nil
}
func (db *Db) ReadListKeys(limit int, keys ...interface{}) ([]string, error) {
	if !work {
		return make([]string, 0), ErrStopped
	}
	db.RWMutex.RLock()
	defer db.RWMutex.RUnlock()
	if len(keys) > len(db.Defkey) {
		return make([]string, 0), ErrWrongParameters
	}
	return db.makeListKeys(limit, keys)

}
func (db *Db) ReadOneRecord() ([]byte, error) {
	if !work {
		return make([]byte, 0), ErrStopped
	}
	db.RWMutex.RLock()
	defer db.RWMutex.RUnlock()
	value, is := db.values[oneKey]
	if !is {
		return make([]byte, 0), ErrKeyNotFound
	}
	return value.Value, nil
}

func (db *Db) ReadRecordFromList(key string) ([]byte, error) {
	if !work {
		return make([]byte, 0), ErrStopped
	}
	db.RWMutex.RLock()
	defer db.RWMutex.RUnlock()
	value, is := db.values[key]
	if !is {
		return make([]byte, 0), ErrKeyNotFound
	}
	return value.Value, nil
}
