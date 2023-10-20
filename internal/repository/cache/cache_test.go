package cache

import (
	"sync"
	"testing"
	"time"

	"github.com/Dsmit05/avost-bot/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestFavouritesCache_AddUser(t *testing.T) {
	tests := []struct {
		name  string
		cache *FavouritesCache
		user  models.User
		want  bool
	}{
		{
			name:  "1 case: want true",
			cache: NewFavouritesCache("test", time.Second),
			user:  models.NewUser(1),
			want:  true,
		},
		{
			name: "2 case: want false",
			cache: &FavouritesCache{
				rw: &sync.RWMutex{},
				m:  map[int64]models.User{1: models.NewUser(1)},
			},
			user: models.NewUser(1),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cache.AddUser(tt.user)
			assert.Equal(t, got, tt.want)
		})
	}
}
