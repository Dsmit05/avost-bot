package markup

import (
	tele "gopkg.in/telebot.v3"
)

func CreateInlineKeyboard(btns ...tele.InlineButton) tele.ReplyMarkup {
	rm := tele.ReplyMarkup{}

	rm.InlineKeyboard = [][]tele.InlineButton{btns}

	return rm
}

// CreateInlineKeyboardSize создает набор кнопок с заданым параметром максимальной ширины.
func CreateInlineKeyboardSize(maxCountWidth int, btns ...tele.InlineButton) tele.ReplyMarkup {
	rm := tele.ReplyMarkup{}

	if maxCountWidth <= 0 {
		return rm
	}

	rows := make([][]tele.InlineButton, (maxCountWidth-1+len(btns))/maxCountWidth)

	for i, b := range btns {
		i /= maxCountWidth
		rows[i] = append(rows[i], b)
	}

	rm.InlineKeyboard = rows

	return rm
}

// CreateInlineKeyboardMany создает наборы кнопок с заданой шириной и максимальным количество в одном сообщении.
func CreateInlineKeyboardMany(maxCountWidth, maxCountBtn int, btns ...tele.InlineButton) []tele.ReplyMarkup {
	if maxCountWidth <= 0 || maxCountBtn <= 0 {
		return []tele.ReplyMarkup{}
	}

	actualSize := (maxCountBtn - 1 + len(btns)) / maxCountBtn

	rms := make([]tele.ReplyMarkup, 0, actualSize)

	for i := 0; i < actualSize; i++ {
		start := i * maxCountBtn
		stop := maxCountBtn + i*maxCountBtn

		if stop > len(btns) {
			stop = len(btns)
		}

		rm := CreateInlineKeyboardSize(maxCountWidth, btns[start:stop]...)

		rms = append(rms, rm)
	}

	return rms
}
