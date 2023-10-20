package text

import (
	"github.com/Dsmit05/avost-bot/internal/bot/core/btn"
	"github.com/Dsmit05/avost-bot/internal/logger"
	"github.com/Dsmit05/avost-bot/pkg/animevost"
	tele "gopkg.in/telebot.v3"
)

type cacheInterface interface {
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

func (i *impl) SearchAnime(c tele.Context) error {
	text := c.Text()

	animes, err := i.aCli.SearchForName(text)
	if err != nil {
		return c.Reply("Ничего не смог найти🙄, попробуйте еще раз.")
	}

	countShowAnime := 5
	if len(animes) < countShowAnime {
		countShowAnime = len(animes)
	}

	for n := 0; n < countShowAnime; n++ {
		imgPre := &tele.Photo{File: tele.FromURL(animes[n].UrlImagePreview)}

		err = c.SendAlbum(tele.Album{imgPre})
		if err != nil {
			logger.Error("b.teleBot.SendAlbum()", err)
		}

		isSub := i.cache.CheckUsersAnime(c.Sender().ID, animes[n].Id)
		btns := btn.PreviewMax(animes[n].Id, isSub, i.mainURL, i.mirrorURL, animes[n].GetPathURl())

		err = c.Send(animes[n].Title, &btns)
		if err != nil {
			logger.Error("SearchAnime() -> c.Send()", err)
		}
	}

	return nil
}
