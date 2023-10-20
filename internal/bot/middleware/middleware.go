package middleware

import (
	"github.com/Dsmit05/avost-bot/internal/logger"
	"github.com/Dsmit05/avost-bot/internal/models"
	"github.com/pkg/errors"
	tele "gopkg.in/telebot.v3"
)

type cacheInterface interface {
	CheckUserRole(id int64, role uint8) bool
}

func New(cache cacheInterface) *impl {
	return &impl{cache: cache}
}

type impl struct {
	cache cacheInterface
}

// Admins check user have rights.
func (i *impl) Admins() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if i.cache.CheckUserRole(c.Sender().ID, models.RoleAdmin) {
				return next(c)
			}

			return errors.New("user not have rights")
		}
	}
}

// LogRecover returns a middleware that recovers a panic happened in handler and log this.
func (i *impl) LogRecover() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			defer func() {
				if r := recover(); r != nil {
					if err, ok := r.(error); ok {
						logger.Error("panic", err)
					} else if s, ok := r.(string); ok {
						logger.Error("panic", errors.New(s))
					} else {
						logger.ErrorAny("panic", r)
					}
				}
			}()

			return next(c)
		}
	}
}
