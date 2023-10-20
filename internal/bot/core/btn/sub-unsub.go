package btn

import (
	"github.com/Dsmit05/avost-bot/internal/bot/core"
	"github.com/Dsmit05/avost-bot/pkg/tele/constructor/btn"
	tele "gopkg.in/telebot.v3"
)

// SubsBtn ★ /fSubs animeID.
func SubsBtn(animeID, prew int) tele.InlineButton {
	return btn.CreateBtnWithArgsInt(core.StarSubs, core.EventSubs, animeID, prew)
}

// UnSubsBtn ✰ /fUnSubs animeID.
func UnSubsBtn(animeID, prew int) tele.InlineButton {
	return btn.CreateBtnWithArgsInt(core.StarUnSubs, core.EventUnSubs, animeID, prew)
}
