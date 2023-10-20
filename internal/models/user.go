package models

const (
	ManageAll     uint8 = 2 // подписан на все обновления
	ManageOnlySub uint8 = 1 // подписан только на избранное
	ManageZero    uint8 = 0 // отписан от всего
)

const (
	RoleDefault uint8 = 0
	RolePro     uint8 = 1

	RoleAdmin uint8 = 5
)

type User struct {
	TelegramID    int64
	SubManageType uint8
	Role          uint8
	SerieInfo     map[int]struct{}
}

func NewUser(telegramID int64) User {
	m := make(map[int]struct{})

	return User{
		TelegramID:    telegramID,
		SubManageType: ManageZero,
		Role:          RoleDefault,
		SerieInfo:     m,
	}
}

func NewUserWithParams(telegramID int64, subManageType uint8, role uint8, serieInfo ...int) User {
	m := make(map[int]struct{}, len(serieInfo))
	for _, v := range serieInfo {
		m[v] = struct{}{}
	}

	return User{
		TelegramID:    telegramID,
		SubManageType: subManageType,
		Role:          role,
		SerieInfo:     m,
	}
}
