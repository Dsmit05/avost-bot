package btn

import (
	"github.com/Dsmit05/avost-bot/internal/bot/core"
	"github.com/Dsmit05/avost-bot/pkg/tele/constructor/btn"
	"github.com/Dsmit05/avost-bot/pkg/tele/constructor/markup"
	tele "gopkg.in/telebot.v3"
)

const (
	TextALl     = "Все обновления"     // 2
	TextOnlySub = "Только избранные"   // 1
	TextZero    = "Отписаться от всех" // 0
)

// SubManagment - Формирует набор кнопок: Все обновления|Только избранные|Отписаться от всех.
func SubManagment() tele.ReplyMarkup {
	btnAll := btn.CreateBtnWithArgsStr(TextALl, core.EventManage, core.ManageAll)             // Все обновления
	btnOnlySub := btn.CreateBtnWithArgsStr(TextOnlySub, core.EventManage, core.ManageOnlySub) // Только избранные
	btnZero := btn.CreateBtnWithArgsStr(TextZero, core.EventManage, core.ManageZero)          // Отписаться от всех

	rm := markup.CreateInlineKeyboard(btnAll, btnOnlySub, btnZero)

	return rm
}

const ManageAllText string = `
Сейчас вы подписаны на все обновления.
Чтобы отписаться нажмите: ` + TextZero + `
Для подписки только на избранное нажмите: ` + TextOnlySub

const ManageOnlySubText string = `
Сейчас вы подписаны только на избранное.
Чтобы отписаться нажмите: ` + TextZero + `
Для подписки на все обновления нажмите: ` + TextALl

const ManageZeroText string = `
Сейчас вы отписаны от всех обновлений.
Чтобы подписаться на все нажмите: ` + TextALl + `
Для подписки только на избранное нажмите: ` + TextOnlySub
