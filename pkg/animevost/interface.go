package animevost

import "github.com/Dsmit05/avost-bot/pkg/animevost/models"

type Interface interface {
	GetPage(page, quantity int) (models.AnimeSpecs, error)
	SearchForName(name string) (models.AnimeSpecs, error)
	GetPlayList(animeID int) (models.Playlists, error)
	GetInfo(animeID int) (models.AnimeSpec, error)
}
