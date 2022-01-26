package data

import "time"

type BaseCtrl struct {
	ID         int       `json:"id"`
	TimeDevice time.Time `json:"dtime"`    // Время устройства
	TechMode   int       `json:"techmode"` //Технологический режим
	/*
		Технологический режим
		1 - выбор ПК по времени по суточной карте ВР-СК;
		2 - выбор ПК по недельной карте ВР-НК;
		3 - выбор ПК по времени по суточной карте, назначенной
		оператором ДУ-СК;
		4 - выбор ПК по недельной карте, назначенной оператором
		ДУ-НК;
		5 - план по запросу оператора ДУ-ПК;
		6 - резервный план (отсутствие точного времени) РП;
		7 – коррекция привязки с ИП;
		8 – коррекция привязки с сервера;
		9 – выбор ПК по годовой карте;
		10 – выбор ПК по ХТ;
		11 – выбор ПК по картограмме;
		12 – противозаторовое управление.
	*/
	Base    bool  `json:"base"` //Если истина то работает по базовой привязке
	PK      int   `json:"pk"`   //Номер плана координации
	CK      int   `json:"ck"`   //Номер суточной карты
	NK      int   `json:"nk"`   //Номер недельной карты
	TMax    int64 `json:"tmax"` //Максимальное время ожидания ответа от сервера в секундах
	TimeOut int64 `json:"tout"` //TimeOut на чтение от контроллера в секундах
}

type Dates struct {
	Name string    `json:"name"`
	Date time.Time `json:"date"`
}

var DatesList = []Dates{
	{Name: "SetupDK", Date: time.UnixMilli(0)},
	{Name: "SetDK", Date: time.UnixMilli(0)},
	{Name: "MounthSets", Date: time.UnixMilli(0)},
	{Name: "WeekSets", Date: time.UnixMilli(0)},
	{Name: "DaySets", Date: time.UnixMilli(0)},
	{Name: "SetCtrl", Date: time.UnixMilli(0)},
	{Name: "SetTimeUse", Date: time.UnixMilli(0)},
	{Name: "TimeDevice", Date: time.UnixMilli(0)},
	{Name: "StatDefine", Date: time.UnixMilli(0)},
	{Name: "PointSet", Date: time.UnixMilli(0)},
	{Name: "UseInput", Date: time.UnixMilli(0)},
}
