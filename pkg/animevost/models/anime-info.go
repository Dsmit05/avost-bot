package models

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/Dsmit05/avost-bot/pkg/animevost/utils"
)

//go:generate easyjson -all
type AnimeSpecs []AnimeSpec

func (a *AnimeSpecs) UnIntersection(newAnimeSpecs AnimeSpecs) AnimeSpecs {
	if a == nil || newAnimeSpecs == nil {
		return nil
	}

	if len(*a) == 0 || len(newAnimeSpecs) == 0 {
		return nil
	}

	hash := make(map[int]struct{}, len(*a))
	for _, oldAnimeSpec := range *a {
		hash[oldAnimeSpec.Id] = struct{}{}
	}

	unInter := make(AnimeSpecs, 0)
	for _, newAnimeSpec := range newAnimeSpecs {
		if _, ok := hash[newAnimeSpec.Id]; !ok {
			unInter = append(unInter, newAnimeSpec)
		}
	}

	return unInter
}

//go:generate easyjson -all
type AnimeSpec struct {
	Id              int         `json:"id"`
	Rating          int         `json:"rating"`
	Votes           int         `json:"votes"`
	Title           string      `json:"title"`
	Description     string      `json:"description"`
	Year            string      `json:"year"`
	UrlImagePreview string      `json:"urlImagePreview"`
	Type            string      `json:"type"`
	Director        string      `json:"director"`
	Series          string      `json:"series"`
	ScreenImage     []string    `json:"screenImage"`
	Timer           interface{} `json:"timer"` // Может быть как Int так и string
}

func (a *AnimeSpec) GetSeries() (map[string]string, error) {
	var m map[string]string
	if err := json.Unmarshal([]byte(strings.Replace(a.Series, "'", "\"", -1)), &m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetScreenImageURLs выводит список ссылок на превью сериала
func (a *AnimeSpec) GetScreenImageURLs() []string {
	if a.ScreenImage == nil || len(a.ScreenImage) == 0 {
		return nil
	}

	imageURLs := make([]string, 0, len(a.ScreenImage))

	for _, v := range a.ScreenImage {
		imageURLs = append(imageURLs, "https://static.openni.ru"+v)
	}

	return imageURLs
}

// GetDescription возвращает отформатированное описание сериала
func (a *AnimeSpec) GetDescription() string {
	r := strings.NewReplacer("<br>", "", "<br />", "\n")
	return r.Replace(a.Description)
}

// GetPreDescription возвращает сокращенное описание сериала
// где symbols это максимальное количество символов
func (a *AnimeSpec) GetPreDescription(symbols int) string {
	r := strings.NewReplacer("<br>", "", "<br />", "\n")
	if len(a.Description) < symbols {
		symbols = len(a.Description)
	}

	preDesk := a.Description[:symbols]

	return r.Replace(preDesk) + "..."
}

func (a *AnimeSpec) GetPathURl() string {
	return strconv.Itoa(a.Id) + "-" + utils.GetTitleNameForURL(a.Title)
}

func (a *AnimeSpec) Info() string {
	var buf strings.Builder

	buf.WriteString("Название: ")
	buf.WriteString(a.Title)
	buf.WriteString("\n\t")

	buf.WriteString("Описание: ")
	buf.WriteString(a.GetDescription())
	buf.WriteString("\n\t")

	buf.WriteString("Рейтинг: ")
	buf.WriteString(strconv.Itoa(a.Rating))
	buf.WriteString("\n\t")

	return buf.String()
}
