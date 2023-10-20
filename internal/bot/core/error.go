package core

import (
	"strings"

	"github.com/pkg/errors"
	tele "gopkg.in/telebot.v3"
)

// CheckErrorForbidden проверяют ошибку 403 от телеграмма
// true - если юзер удалил бота или сам исчез.
func CheckErrorForbidden(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, tele.ErrBlockedByUser) {
		return true
	}

	if errors.Is(err, tele.ErrNotStartedByUser) {
		return true
	}

	if errors.Is(err, tele.ErrUserIsDeactivated) {
		return true
	}

	if strings.Contains(err.Error(), "(403)") {
		return true
	}

	return false
}
