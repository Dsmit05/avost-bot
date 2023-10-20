package callback

import (
	"strconv"

	"github.com/Dsmit05/avost-bot/internal/bot/core"
	"github.com/Dsmit05/avost-bot/internal/bot/core/btn"
	"github.com/Dsmit05/avost-bot/internal/logger"
	"github.com/Dsmit05/avost-bot/pkg/animevost"
	telebtn "github.com/Dsmit05/avost-bot/pkg/tele/constructor/btn"
	"github.com/pkg/errors"
	tele "gopkg.in/telebot.v3"
)

type cacheInterface interface {
	UpdateUserSubManageType(userID int64, newSubManageType uint8) bool
	AddAnime(userID int64, animeID int) bool
	DelAnime(userID int64, animeID int) bool
	CheckUsersAnime(userID int64, animeID int) bool
}

func New(aCli animevost.Interface, cache cacheInterface, mainURL, mirrorURL string) *impl {
	return &impl{aCli: aCli, cache: cache, mainURL: mainURL, mirrorURL: mirrorURL}
}

type impl struct {
	aCli               animevost.Interface
	cache              cacheInterface
	mainURL, mirrorURL string
}

func (i *impl) Call(c tele.Context) error {
	event, args := telebtn.GetBtnArgsStr(c.Data())
	if len(args) == 0 {
		return errors.New("not have args")
	}

	switch event {
	case core.EventInfo:
		return i.sendSeriesInfo(c, args)
	case core.EventSubs:
		return i.sendEventSubs(c, args)
	case core.EventUnSubs:
		return i.sendEventUnSubs(c, args)
	case core.EventManage:
		return i.sendSubManagment(c, args)
	}

	return errors.New("not have event")
}

func (i *impl) sendEventSubs(c tele.Context, args []string) error {
	intV, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.Wrap(err, "strconv.Atoi() error:")
	}

	if !i.cache.DelAnime(c.Sender().ID, intV) {
		return c.Reply("уже удален")
	}

	prew, err := strconv.Atoi(args[1])
	if err != nil {
		return errors.Wrap(err, "strconv.Atoi() error:")
	}

	spec, err := i.aCli.GetInfo(intV)
	if err != nil {
		return errors.Wrap(err, "cli.GetInfo() error")
	}

	if prew == 0 {
		btns := btn.PreviewShort(intV, false, i.mainURL, i.mirrorURL, spec.GetPathURl())

		return c.Edit(&btns)
	} else {
		btns := btn.PreviewMax(intV, false, i.mainURL, i.mirrorURL, spec.GetPathURl())

		return c.Edit(&btns)
	}
}

func (i *impl) sendEventUnSubs(c tele.Context, args []string) error {
	intV, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.Wrap(err, "strconv.Atoi() error:")
	}

	if !i.cache.AddAnime(c.Sender().ID, intV) {
		return c.Reply("уже удален")
	}

	prew, err := strconv.Atoi(args[1])
	if err != nil {
		return errors.Wrap(err, "strconv.Atoi() error:")
	}

	spec, err := i.aCli.GetInfo(intV)
	if err != nil {
		return errors.Wrap(err, "cli.GetInfo() error")
	}

	if prew == 0 {
		btns := btn.PreviewShort(intV, true, i.mainURL, i.mirrorURL, spec.GetPathURl())

		return c.Edit(&btns)
	} else {
		btns := btn.PreviewMax(intV, true, i.mainURL, i.mirrorURL, spec.GetPathURl())

		return c.Edit(&btns)
	}
}

func (i *impl) sendSubManagment(c tele.Context, args []string) error {
	btns := btn.SubManagment()

	switch args[0] {
	case core.ManageAll:
		i.cache.UpdateUserSubManageType(c.Sender().ID, 2)
		return c.Edit(btn.ManageAllText, &btns)
	case core.ManageOnlySub:
		i.cache.UpdateUserSubManageType(c.Sender().ID, 1)
		return c.Edit(btn.ManageOnlySubText, &btns)
	case core.ManageZero:
		i.cache.UpdateUserSubManageType(c.Sender().ID, 0)
		return c.Edit(btn.ManageZeroText, &btns)
	}

	return nil
}

// sendSeriesInfo выводит подробное описание аниме.
func (i *impl) sendSeriesInfo(c tele.Context, args []string) error {
	animeID, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.Wrap(err, "strconv.Atoi() error:")
	}

	spec, err := i.aCli.GetInfo(animeID)
	if err != nil {
		return errors.Wrap(err, "cli.GetInfo() error")
	}

	// Получаем список превью картинок и отправляем
	urls := spec.GetScreenImageURLs()
	if urls != nil {
		imgs := make(tele.Album, 0, len(urls))

		for _, v := range urls {
			imgs = append(imgs, &tele.Photo{File: tele.FromURL(v)})
		}

		err = c.SendAlbum(imgs)
		if err != nil {
			logger.Error("SendSeriesInfo -> c.SendAlbum() error", err)
		}
	}

	// Формируем кнопки под картинками
	isSub := i.cache.CheckUsersAnime(c.Sender().ID, spec.Id)
	btns := btn.PreviewShort(animeID, isSub, i.mainURL, i.mirrorURL, spec.GetPathURl())

	err = c.Send(spec.Info(), &btns)
	if err != nil {
		return errors.Wrap(err, "c.Send() error")
	}

	// Получаем список серий и отправляем
	series, err := spec.GetSeries()
	if err != nil {
		return errors.Wrap(err, "spec.GetSeries() error")
	}

	markupSeries := btn.Series(series, 5, 100)
	for _, mk := range markupSeries {
		err = c.Send("Список серий для просмотра:", &mk)
		if err != nil {
			return errors.Wrap(err, "sendSeriesInfo error:")
		}
	}

	return nil
}
