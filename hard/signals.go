package hard

type bits struct {
	word int //номер слова
	bit  int //номер бита
}
type DI struct {
	nomer    int  //Номер входа
	value    bits //Место хранения значения
	counter  int  //Место хранения счетчика
	state    int  //Место хранения/установки режима работы
	reset    bits //Место сброса состояний в режимах Вход с удержанием Тригер и Счетчик
	bounce   int  //Место хранения времени антидребезга
	blocked  int  //Место хранения времени удержания входа
	bl_state bits //Место удерживаемых состояний
	front    bits //Место хранения активных фронтов
}
type AI struct {
	nomer int //Номер входа
	value int //Место хранения значения

}

type DO struct {
	nomer int  //Номер входа
	value bits //Место хранения значения

}
