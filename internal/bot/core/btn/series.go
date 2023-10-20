package btn

import (
	"regexp"
	"sort"
	"strconv"

	"github.com/Dsmit05/avost-bot/pkg/tele/constructor/btn"
	"github.com/Dsmit05/avost-bot/pkg/tele/constructor/markup"
	tele "gopkg.in/telebot.v3"
)

var regInt = regexp.MustCompile(`\d+`)

func Series(m map[string]string, maxCountWidth, maxCountBtn int) []tele.ReplyMarkup {
	btns := make([]tele.InlineButton, 0)

	for k, v := range m {
		urlBtn := btn.CreateBtnURL(k, "http://video.aniland.org/"+v+".mp4")
		btns = append(btns, urlBtn)
	}

	sort.Slice(btns, func(i, j int) bool {
		first := btns[i].Text
		second := btns[j].Text

		fInt, _ := strconv.Atoi(regInt.FindString(first))
		lInt, _ := strconv.Atoi(regInt.FindString(second))
		return fInt < lInt
	})

	mk := markup.CreateInlineKeyboardMany(maxCountWidth, maxCountBtn, btns...)

	return mk
}
