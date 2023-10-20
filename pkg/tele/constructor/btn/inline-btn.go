package btn

import (
	tele "gopkg.in/telebot.v3"
)

// CreateBtnWithArgsStr создает уникальную кнопку со следующей сигнатурой: /f<<Event>> <<Arg>> <<Arg>>...
func CreateBtnWithArgsStr(text, event string, args ...string) tele.InlineButton {
	unique := EncodeStr(event, args...)

	return tele.InlineButton{Text: text, Unique: unique}
}

// CreateBtnWithArgsInt  /f<<Event>> <<Int>> <<int>>...
func CreateBtnWithArgsInt(text, event string, args ...int) tele.InlineButton {
	unique := EncodeInt(event, args...)

	return tele.InlineButton{Text: text, Unique: unique}
}

func GetBtnArgsStr(unique string) (event string, args []string) {
	if len(unique) < 2 {
		return
	}

	event, args = DecodeStr(unique[1:])

	return
}

// CreateBtnURL кнопка с ссылкой
func CreateBtnURL(text, url string) tele.InlineButton {
	return tele.InlineButton{Text: text, URL: url}
}
