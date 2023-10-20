package cache

import (
	"context"
	"encoding/gob"
	"os"
	"sync"
	"time"

	"github.com/Dsmit05/avost-bot/internal/logger"
	"github.com/Dsmit05/avost-bot/internal/models"
	"github.com/pkg/errors"
)

type FavouritesCache struct {
	path   string
	rw     *sync.RWMutex
	m      map[int64]models.User // user id text ids
	ticker *time.Ticker
}

func NewFavouritesCache(path string, timeSave time.Duration) *FavouritesCache {
	ticker := time.NewTicker(timeSave)

	return &FavouritesCache{
		path:   path,
		rw:     new(sync.RWMutex),
		m:      make(map[int64]models.User),
		ticker: ticker,
	}
}

func (f *FavouritesCache) Start(ctx context.Context) error {
	if err := f.Load(); err != nil {
		return err
	}

	for {
		select {
		case <-f.ticker.C:
			if err := f.Save(); err != nil {
				logger.Error("cache Save", err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}

}

func (f *FavouritesCache) Stop() error {
	f.ticker.Stop()

	return f.Save()
}

// AddUser добавляет пользователя в кеш, если уже есть возвращает false.
func (f *FavouritesCache) AddUser(user models.User) bool {
	f.rw.RLock()
	_, ok := f.m[user.TelegramID]
	if ok {
		f.rw.RUnlock()
		return false
	}
	f.rw.RUnlock()

	f.rw.Lock()
	f.m[user.TelegramID] = user
	f.rw.Unlock()

	return true
}

func (f *FavouritesCache) GetUser(userID int64) (models.User, bool) {
	f.rw.RLock()
	defer f.rw.RUnlock()
	user, ok := f.m[userID]

	return user, ok
}

// DelUser удаляет пользователя.
func (f *FavouritesCache) DelUser(userID int64) bool {
	f.rw.Lock()
	defer f.rw.Unlock()

	_, ok := f.m[userID]
	if !ok {
		return false
	}

	delete(f.m, userID)

	return true
}

// UpdateUserRole обновляет роль пользователя.
func (f *FavouritesCache) UpdateUserRole(userID int64, newRole uint8) bool {
	f.rw.RLock()
	user, ok := f.m[userID]
	if !ok {
		f.rw.RUnlock()
		return false
	}
	f.rw.RUnlock()

	user.Role = newRole

	f.rw.Lock()
	f.m[userID] = user
	f.rw.Unlock()

	return true
}

func (f *FavouritesCache) UpdateUserSubManageType(userID int64, newSubManageType uint8) bool {
	f.rw.RLock()
	user, ok := f.m[userID]
	if !ok {
		f.rw.RUnlock()
		return false
	}
	f.rw.RUnlock()

	user.SubManageType = newSubManageType

	f.rw.Lock()
	f.m[userID] = user
	f.rw.Unlock()

	return true
}

func (f *FavouritesCache) GetUserSubManageType(userID int64) uint8 {
	f.rw.RLock()
	user, ok := f.m[userID]
	if !ok {
		f.rw.RUnlock()
		return 0
	}
	f.rw.RUnlock()

	return user.SubManageType
}

// AddAnime добавляет пользователю аниме.
func (f *FavouritesCache) AddAnime(userID int64, animeID int) bool {
	f.rw.Lock()
	defer f.rw.Unlock()

	_, ok := f.m[userID]
	if !ok {
		f.m[userID] = models.NewUser(userID)
	}

	if _, animeOK := f.m[userID].SerieInfo[animeID]; animeOK {
		return false
	}

	f.m[userID].SerieInfo[animeID] = struct{}{}

	return true
}

func (f *FavouritesCache) DelAnime(userID int64, animeID int) bool {
	f.rw.Lock()
	defer f.rw.Unlock()

	_, ok := f.m[userID]
	if !ok {
		return false
	}

	_, animeOk := f.m[userID].SerieInfo[animeID]
	if !animeOk {
		return false
	}

	delete(f.m[userID].SerieInfo, animeID)

	return true
}

func (f *FavouritesCache) GetAllUserAnime(userID int64) ([]int, bool) {
	f.rw.RLock()
	user, ok := f.m[userID]
	if !ok {
		f.rw.RUnlock()
		return nil, false
	}
	f.rw.RUnlock()

	if user.SerieInfo == nil || len(user.SerieInfo) == 0 {
		return nil, false
	}

	sl := make([]int, 0, len(user.SerieInfo))
	for k, _ := range user.SerieInfo {
		sl = append(sl, k)
	}

	return sl, true
}

func (f *FavouritesCache) GetAllMap() map[int64]models.User {
	f.rw.RLock()
	defer f.rw.RUnlock()

	return f.m
}

func (f *FavouritesCache) CheckUserRole(id int64, role uint8) bool {
	f.rw.RLock()
	defer f.rw.RUnlock()

	user, ok := f.m[id]
	if !ok {
		return false
	}

	if user.Role == role {
		return true
	}

	return false
}

// CheckUsersAnime проверяет есть ли у пользователя данное аниме
func (f *FavouritesCache) CheckUsersAnime(userID int64, animeID int) bool {
	f.rw.RLock()
	defer f.rw.RUnlock()
	user, ok := f.m[userID]
	if !ok {
		return false
	}

	if user.SerieInfo == nil {
		return false
	}

	_, ok = user.SerieInfo[animeID]
	if !ok {
		return false
	}

	return true
}

// GetUsersID возвращает список  ID пользователей.
func (f *FavouritesCache) GetUsersID() ([]int64, bool) {
	f.rw.RLock()
	defer f.rw.RUnlock()

	if len(f.m) == 0 {
		return nil, false
	}

	sl := make([]int64, 0, len(f.m))

	for k, _ := range f.m {
		sl = append(sl, k)
	}

	return sl, true
}

// GetUsers возвращает список пользователей.
func (f *FavouritesCache) GetUsers() ([]models.User, bool) {
	f.rw.RLock()
	defer f.rw.RUnlock()

	if len(f.m) == 0 {
		return nil, false
	}

	sl := make([]models.User, 0, len(f.m))

	for _, v := range f.m {
		sl = append(sl, v)
	}

	return sl, true
}

// GetUsersAdmins возвращает список админов.
func (f *FavouritesCache) GetUsersAdmins() []models.User {
	f.rw.RLock()
	defer f.rw.RUnlock()

	sl := make([]models.User, 0)
	for _, v := range f.m {
		if v.Role == models.RoleAdmin {
			sl = append(sl, v)
		}
	}

	return sl
}

func (f *FavouritesCache) Save() error {
	fp, err := os.Create(f.path)
	if err != nil {
		return errors.Wrapf(err, "os.Create(%s) error", f.path)
	}
	defer fp.Close()

	enc := gob.NewEncoder(fp)

	f.rw.RLock()
	err = enc.Encode(&f.m)
	f.rw.RUnlock()

	return errors.Wrap(err, "enc.Encode() error")
}

func (f *FavouritesCache) Load() error {
	fp, err := os.Open(f.path)
	if err != nil {
		return errors.Wrapf(err, "os.Open(%s) error", f.path)
	}

	defer fp.Close()

	var m map[int64]models.User
	dec := gob.NewDecoder(fp)

	err = dec.Decode(&m)
	if err != nil {
		return errors.Wrap(err, "dec.Decode() error")
	}

	f.rw.Lock()
	f.m = m
	f.rw.Unlock()

	return nil
}

func (f *FavouritesCache) GetCountUsers() int {
	f.rw.RLock()
	defer f.rw.RUnlock()
	return len(f.m)
}
