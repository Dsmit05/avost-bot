package btn

import (
	"github.com/Dsmit05/avost-bot/internal/bot/core"
	"github.com/Dsmit05/avost-bot/pkg/tele/constructor/btn"
	"github.com/Dsmit05/avost-bot/pkg/tele/constructor/markup"
	tele "gopkg.in/telebot.v3"
)

const (
	TextDetailed  = "Подробнее"
	TextMainURL   = "Оф сайт"
	TextMirrorURL = "Зеркало"
)

// PreviewMax - Формирует набор кнопок: Подробнее|Оф сайт|Зеркало|★.
func PreviewMax(animeID int, isSub bool, mainURL, mirrorURL, path string) tele.ReplyMarkup {
	btnMoreDetailed := btn.CreateBtnWithArgsInt(TextDetailed, core.EventInfo, animeID) // Подробнее
	btnMainURL := btn.CreateBtnURL(TextMainURL, CreateMainURL(mainURL, path))          // Оф сайт
	btnMirrorURL := btn.CreateBtnURL(TextMirrorURL, CreateMirrorURL(mirrorURL, path))  // Зекрало
	var btnSubsOrUnSubs tele.InlineButton                                              // ★

	if isSub {
		btnSubsOrUnSubs = SubsBtn(animeID, 1)
	} else {
		btnSubsOrUnSubs = UnSubsBtn(animeID, 1)
	}

	rm := markup.CreateInlineKeyboard(btnMoreDetailed, btnMainURL, btnMirrorURL, btnSubsOrUnSubs)

	return rm
}

// PreviewShort - Формирует набор кнопок: Оф сайт|Зеркало|★.
func PreviewShort(animeID int, isSub bool, mainURL, mirrorURL, path string) tele.ReplyMarkup {
	btnMainURL := btn.CreateBtnURL(TextMainURL, CreateMainURL(mainURL, path))         // Оф сайт
	btnMirrorURL := btn.CreateBtnURL(TextMirrorURL, CreateMirrorURL(mirrorURL, path)) // Зекрало
	var btnSubsOrUnSubs tele.InlineButton                                             // ★

	if isSub {
		btnSubsOrUnSubs = SubsBtn(animeID, 0)
	} else {
		btnSubsOrUnSubs = UnSubsBtn(animeID, 0)
	}

	rm := markup.CreateInlineKeyboard(btnMainURL, btnMirrorURL, btnSubsOrUnSubs)

	return rm
}

func CreateMainURL(mainURL, path string) string {
	return mainURL + "/tip/tv/" + path + ".html"
}

func CreateMirrorURL(mirrorURL, path string) string {
	return mirrorURL + "/tip/tv/" + path + ".html"
}
