package command

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/Dsmit05/avost-bot/internal/bot/controllers"
	"github.com/Dsmit05/avost-bot/internal/bot/core/btn"
	"github.com/Dsmit05/avost-bot/internal/logger"
	"github.com/Dsmit05/avost-bot/internal/models"
	"github.com/Dsmit05/avost-bot/pkg/animevost"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

type cacheInterface interface {
	AddUser(user models.User) bool
	GetAllUserAnime(userID int64) ([]int, bool)
	CheckUsersAnime(userID int64, animeID int) bool
	GetUsersAdmins() []models.User
	GetUserSubManageType(userID int64) uint8
	GetCountUsers() int
}

func New(aCli animevost.Interface, cache cacheInterface, mainURL, mirrorURL string) *impl {
	return &impl{aCli: aCli, cache: cache, mainURL: mainURL, mirrorURL: mirrorURL}
}

type impl struct {
	aCli               animevost.Interface
	cache              cacheInterface
	mainURL, mirrorURL string
}

func (i *impl) Start(c tele.Context) error {
	i.cache.AddUser(models.NewUser(c.Sender().ID))

	return c.Send(controllers.StartText)
}

func (i *impl) Help(c tele.Context) error {
	return c.Send(controllers.HelpText)
}

// Favourites /favourites - избранное.
func (i *impl) Favourites(c tele.Context) error {
	animeIDs, ok := i.cache.GetAllUserAnime(c.Sender().ID)
	if !ok {
		return c.Reply("Вы пока ничего не добавили")
	}

	for _, animeID := range animeIDs {
		anime, err := i.aCli.GetInfo(animeID)
		if err != nil {
			logger.Error("i.aCli.GetInfo(animeID)", err)

			continue
		}

		btns := btn.PreviewMax(animeID, true, i.mainURL, i.mirrorURL, anime.GetPathURl())

		err = c.Send(anime.Title, &btns)
		if err != nil {
			logger.Error("Favourites() -> c.Send()", err)
		}
	}

	return nil
}

// Last /last5 - последние вышедшие 5 аниме.
func (i *impl) Last(c tele.Context) error {
	animes, err := i.aCli.GetPage(1, 5)
	if err != nil {
		return errors.Wrap(err, " i.aCli.GetPage(1, 5) error")
	}

	userID := c.Sender().ID

	for _, animeSpec := range animes {
		// Отправка фото из UrlImagePreview
		imgPre := &tele.Photo{File: tele.FromURL(animeSpec.UrlImagePreview)}

		err = c.SendAlbum(tele.Album{imgPre})
		if err != nil {
			logger.Error("b.teleBot.SendAlbum()", err)
		}

		// Отправка текста и кнопок для пользователя
		isSub := i.cache.CheckUsersAnime(userID, animeSpec.Id)
		btns := btn.PreviewMax(animeSpec.Id, isSub, i.mainURL, i.mirrorURL, animeSpec.GetPathURl())

		text := animeSpec.Title + "\n" + animeSpec.GetPreDescription(200)

		err = c.Send(text, &btns)
		if err != nil {
			logger.Error("c.Send()", err)
		}
	}

	return nil
}

// Random /random - случайное аниме.
func (i *impl) Random(c tele.Context) error {
	n := 1 + rand.Intn(2850)

	animes, err := i.aCli.GetPage(n, 1)
	if err != nil {
		return errors.Wrap(err, "i.aCli.GetPage(n, 1) error")
	}

	if len(animes) == 0 {
		return errors.New("not have anime")
	}

	anime := animes[0]

	imgPre := &tele.Photo{File: tele.FromURL(anime.UrlImagePreview)}

	err = c.SendAlbum(tele.Album{imgPre})
	if err != nil {
		return errors.Wrap(err, "SendAlbum() error")
	}

	isSub := i.cache.CheckUsersAnime(c.Sender().ID, anime.Id)
	btns := btn.PreviewMax(anime.Id, isSub, i.mainURL, i.mirrorURL, anime.GetPathURl())

	return c.Send(anime.Title+"\n"+anime.GetPreDescription(250), &btns)
}

// Feedback /fb текст сообщения.
func (i *impl) Feedback(c tele.Context) error {
	msg := c.Data()
	if len(msg) <= 5 {
		return c.Reply("Добавьте после команды сообщение")
	}

	c.Reply("Спасибо за фидбек")

	user := c.Sender()
	if user == nil {
		return errors.New("user is nil")
	}

	msgID := c.Message().ID

	text := fmt.Sprintf("ID: %d, Nik: %s, Name: %s %s\nmsgID: %d: %s", user.ID, user.Username, user.FirstName, user.LastName, msgID, msg)

	for _, admin := range i.cache.GetUsersAdmins() {
		chat, err := c.Bot().ChatByID(admin.TelegramID)
		if err != nil {
			logger.Errorf("Feedback ->  c.Bot().ChatByUsername(admin) error", zap.String("error", err.Error()))

			continue
		}

		_, err = c.Bot().Send(chat, text)
		if err != nil {
			logger.Errorf("Feedback ->  c.Bot().Send(chatID, text) error", zap.String("error", err.Error()))
		}
	}

	return nil
}

// Subscription /sub - управление подписками.
func (i *impl) Subscription(c tele.Context) error {
	var text string

	manageType := i.cache.GetUserSubManageType(c.Sender().ID)
	switch manageType {
	case models.ManageAll:
		text = btn.ManageAllText
	case models.ManageOnlySub:
		text = btn.ManageOnlySubText
	case models.ManageZero:
		text = btn.ManageZeroText
	}

	btns := btn.SubManagment()

	return c.Send(text, &btns)
}

// Len /len - количество пользователей.
func (i *impl) Len(c tele.Context) error {
	counts := i.cache.GetCountUsers()

	return c.Send(strconv.Itoa(counts))
}

func (i *impl) ID(c tele.Context) error {
	id := strconv.FormatInt(c.Sender().ID, 10)
	return c.Send(id)
}
