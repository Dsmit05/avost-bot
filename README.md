# avost-bot
https://t.me/AnimeVostBot

### About

Сервис для работы с сайтом animevost.
Умеет искать, сохранять и просмотривать серии аниме.

Для общения c пользователями и администрирования реализует апи(:8080) + swagger.

Бот хранит состояние пользователей, логи и статистку в папке /data

**Main command**

```
/favourites - Избранное   
/last5 - 5 последних аниме  
/random - Случайное аниме  
/sub - Подписка на обновления  
/help - Помощь
/fb - Отзыв
```

**Admin command**

```
/len - Количество пользователей   
/stat - Статистика за сутки по часам
```

### Config

Конфигурация задается через переменные окружения

| Название   | Описание                      |
|------------|-------------------------------|
| BOT_TOKEN  | Токен бота выданный BotFather |
| MAIN_URL   | Основной адресс сайта         |
| MIRROR_URL | Зеркало                       |
| JWT        | Секретный ключ для JWT        |
| ADDRESS    | Адресс api бота               |

### Role model

В зависимости от роли пользователя предоставляет дополнительные возможности.

Ограничения по ролям, на данный момент используется только default и admin.

| Название | Тип         | Кеш |
|----------|-------------|-----|
| default  | RoleDefault | 0   |
| pro      | RolePro     | 1   |
| admin    | RoleAdmin   | 5   |

Тип подписок на обновления.

| Название | Тип            | Кеш | Описание            |
|----------|----------------|-----|---------------------|
| Zero     | ManageZero     | 0   | Отписаться от всех  |
| OnlySub  | ManageOnlySub  | 1   | Только избранные    |
| All      | ManageAll      | 2   | Все обновления      |

### Deploy
Сервис разворачивается на виртуальную машину через гитлаб ci/cd.
