package bot

import (
	"github.com/Dsmit05/avost-bot/pkg/stat"
	"runtime/debug"
	"time"

	"github.com/Dsmit05/avost-bot/internal/bot/controllers/callback"
	"github.com/Dsmit05/avost-bot/internal/bot/controllers/command"
	"github.com/Dsmit05/avost-bot/internal/bot/controllers/text"
	"github.com/Dsmit05/avost-bot/internal/bot/core"
	"github.com/Dsmit05/avost-bot/internal/bot/core/btn"
	"github.com/Dsmit05/avost-bot/internal/bot/middleware"
	"github.com/Dsmit05/avost-bot/internal/config"
	"github.com/Dsmit05/avost-bot/internal/logger"
	"github.com/Dsmit05/avost-bot/internal/models"
	"github.com/Dsmit05/avost-bot/internal/repository"
	"github.com/Dsmit05/avost-bot/pkg/animevost"
	avmodels "github.com/Dsmit05/avost-bot/pkg/animevost/models"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

type Bot struct {
	teleBot  *tele.Bot
	cfg      *config.AnimeBotConfg
	avClient animevost.Interface
	cache    repository.CacheInterface
	stat     *stat.Statistic
}

func NewBot(cfg *config.AnimeBotConfg, avClient animevost.Interface, cache repository.CacheInterface, stat *stat.Statistic) (*Bot, error) {
	settings := tele.Settings{
		Token:   cfg.Token,
		Poller:  &tele.LongPoller{Timeout: cfg.LongPollerTimeout},
		OnError: LogError,
	}

	b, err := tele.NewBot(settings)
	if err != nil {
		return nil, errors.Wrap(err, "tele.NewBot() error")
	}

	return &Bot{teleBot: b, cfg: cfg, avClient: avClient, cache: cache, stat: stat}, nil
}

func (b *Bot) InitControllersV1() *Bot {
	commandController := command.New(b.avClient, b.cache, b.cfg.MainURL, b.cfg.MirrorURL)
	textController := text.New(b.avClient, b.cache, b.cfg.MainURL, b.cfg.MirrorURL)
	callbackController := callback.New(b.avClient, b.cache, b.cfg.MainURL, b.cfg.MirrorURL)
	middlewareController := middleware.New(b.cache)

	mainGroup := b.teleBot.Group()
	mainGroup.Use(middlewareController.LogRecover(), b.stat.BotMiddleware())
	{
		mainGroup.Handle("/start", commandController.Start)
		mainGroup.Handle("/help", commandController.Help)
		mainGroup.Handle("/last5", commandController.Last)
		mainGroup.Handle("/favourites", commandController.Favourites)
		mainGroup.Handle("/random", commandController.Random)
		mainGroup.Handle("/fb", commandController.Feedback)
		mainGroup.Handle("/sub", commandController.Subscription)
		mainGroup.Handle("/len", commandController.Len, middlewareController.Admins())
		mainGroup.Handle("/id", commandController.ID)
		mainGroup.Handle("/stat", b.stat.HandlerStat, middlewareController.Admins())

		// text
		mainGroup.Handle(tele.OnText, textController.SearchAnime)

		// callback
		mainGroup.Handle(tele.OnCallback, callbackController.Call)

	}

	return b
}

func (b *Bot) Start() {
	go b.StartSendUpdate()
	b.teleBot.Start()
}

func (b *Bot) Stop() {
	b.teleBot.Stop()
}

func (b *Bot) Replay(msgID int, tgID int64, text string) error {
	chat := tele.Chat{ID: tgID}
	msg := tele.Message{
		ID:   msgID,
		Chat: &chat,
	}
	_, err := b.teleBot.Reply(&msg, text)

	return err
}

func (b *Bot) Send(tgID int64, text string) error {
	chat := tele.Chat{ID: tgID}

	_, err := b.teleBot.Send(&chat, text)

	return err
}

func (b *Bot) SendAll(text string) error {
	users, ok := b.cache.GetUsersID()
	if !ok {
		return errors.New("not have user")
	}

	var err error

	for _, id := range users {
		chat := tele.Chat{ID: id}
		_, errBot := b.teleBot.Send(&chat, text)
		if core.CheckErrorForbidden(errBot) {
			logger.Infof("SendAll -> delete user", zap.Int64("userID", id))
			b.cache.DelUser(id)
		}

		err = multierr.Append(err, errBot)
	}

	return err
}

func (b *Bot) GetUserInfo(id int64) (*models.UserFullInfo, error) {
	cacheInfo, ok := b.cache.GetUser(id)
	if !ok {
		return nil, errors.New("not have user in cache")
	}

	chat, err := b.teleBot.ChatByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "b.teleBot.ChatByID()")
	}

	var user models.UserFullInfo

	user.ID = id
	user.Username = chat.Username
	user.FirstName = chat.FirstName
	user.LastName = chat.LastName

	switch cacheInfo.SubManageType {
	case models.ManageAll:
		user.SubManageType = core.ManageAll
	case models.ManageOnlySub:
		user.SubManageType = core.ManageOnlySub
	case models.ManageZero:
		user.SubManageType = core.ManageZero
	default:
		user.SubManageType = "unknown"
	}

	switch cacheInfo.Role {
	case models.RoleDefault:
		user.Role = "default"
	case models.RolePro:
		user.Role = "pro"
	case models.RoleAdmin:
		user.Role = "admin"
	default:
		user.SubManageType = "unknown"
	}

	if len(cacheInfo.SerieInfo) == 0 {
		return &user, nil
	}

	sl := make([]models.Anime, 0, len(cacheInfo.SerieInfo))
	for k, _ := range cacheInfo.SerieInfo {
		anime, _ := b.avClient.GetInfo(k)

		sl = append(sl, models.Anime{
			Name: anime.Title,
			URL:  btn.CreateMirrorURL(b.cfg.MirrorURL, anime.GetPathURl()),
		})
	}

	user.Favorites = sl

	return &user, nil
}

// StartSendUpdate проверяет обновления на сайте и если появились новые серии отправляет их всем юзерам.
func (b *Bot) StartSendUpdate() {
	logger.Info("StartSendUpdate()", "start")

	defer func() {
		if r := recover(); r != nil {
			logger.Error("panic", errors.New(string(debug.Stack())))

			go b.StartSendUpdate()
		}

	}()

	var updates avmodels.AnimeSpecs // Содержит последние вышедшие аниме

	animeSpecs, errGetPage := b.avClient.GetPage(1, 5)
	if errGetPage != nil {
		logger.Error("b.avClient.GetPage()", errGetPage)
	}

	updates = animeSpecs

	for {
		time.Sleep(b.cfg.LastAnimeUpdateTimeout)

		animes, err := b.avClient.GetPage(1, 5)
		if err != nil {
			logger.Error("b.avClient.GetPage()", err)
			continue
		}

		// При первом запуске может не пройти запрос, поэтому при повторном переопределяем
		if updates == nil {
			updates = animes
			continue
		}

		unInter := updates.UnIntersection(animes)
		updates = animes
		if unInter == nil || len(unInter) == 0 {
			continue
		}

		users, ok := b.cache.GetUsers()
		if !ok {
			logger.Error("b.cache.GetUsers()", errors.New("not have user in cache"))
			continue
		}

		logger.Info("StartSendUpdate", "Появилось новое аниме")

		b.SendUsersAnime(unInter, users...)
	}
}

func (b *Bot) SendUsersAnime(animeSpecs avmodels.AnimeSpecs, users ...models.User) {
	for _, animeSpec := range animeSpecs {
		for _, user := range users {
			if user.SubManageType == models.ManageZero {
				continue
			}

			if user.SubManageType == models.ManageOnlySub {
				if _, ok := user.SerieInfo[animeSpec.Id]; !ok {
					continue
				}
			}

			chatID := tele.ChatID(user.TelegramID)

			// Отправка фото из UrlImagePreview
			imgPre := &tele.Photo{File: tele.FromURL(animeSpec.UrlImagePreview)}
			_, err := b.teleBot.SendAlbum(chatID, tele.Album{imgPre})
			if err != nil {
				logger.Error("SendUsersAnime -> b.teleBot.SendAlbum()", err)

				if core.CheckErrorForbidden(err) {
					logger.Infof("SendUsersAnime -> delete user", zap.Int64("userID", user.TelegramID))

					if !b.cache.DelUser(user.TelegramID) {
						logger.Errorf("SendUsersAnime -> b.cache.DelUser() false", zap.Int64("userID", user.TelegramID))
					}
				}

				continue
			}

			// Отправка текста и кнопок для пользователя
			isSub := b.cache.CheckUsersAnime(user.TelegramID, animeSpec.Id)
			btns := btn.PreviewMax(animeSpec.Id, isSub, b.cfg.MainURL, b.cfg.MirrorURL, animeSpec.GetPathURl())

			textData := "Новая серия:\n" + animeSpec.Title + "\n" + animeSpec.GetPreDescription(200)
			_, err = b.teleBot.Send(chatID, textData, &btns)
			if err != nil {
				logger.Error("SendUsersAnime -> b.teleBot.Send()", err)

				if core.CheckErrorForbidden(err) {
					logger.Infof("SendUsersAnime -> delete user", zap.Int64("userID", user.TelegramID))

					if !b.cache.DelUser(user.TelegramID) {
						logger.Errorf("SendUsersAnime -> b.cache.DelUser() false", zap.Int64("userID", user.TelegramID))
					}
				}

				continue
			}

		}
	}
}
