# GophKeeper
Менеджер паролей GophKeeper

## Сервер

### [API спецификация](docs/api/swagger.json)
У сервиса есть с актуальная API спецификация. Ее можно использовать в Swagger UI или Postman для импорта.

### Эксплуатация

#### Логи
Пишутся в формате json. Пример лога http-запроса:

```json
{"content_length":0,"duration_ms":84,"level":"info","method":"POST","msg":"Request info","status":200,"time":"2024-03-05T17:50:57+03:00","uri":"/api/v1/auth/register"}
```

#### Хелсчек
Ендпоинт `/health` отдает 200 http-код.

#### Конфигурация

Переменные необходимые для запуска приложения:

| Параметр            | Описание                         | По умолчанию                                       |
|---------------------|----------------------------------|----------------------------------------------------|
| ADDR                | Хост и порт для запуска          | 0.0.0.0:8080                                       | 
| DB_TIMEOUT          | Таймаут операций БД              | 15s                                                |
| JWT_EXPIRATION      | Время жизни JWT                  | 600s                                               |
| RAW_JWK             | JSON Web Keys                    | My secret keys                                     |
| POSTGRES_URL        | Postgres URL                     | postgresql://admin:admin@localhost:5432/gophkeeper |
| CRYPTO_SECRET       | Приватный ключ шифрования данных | 1234567812345678                                   |
| READ_HEADER_TIMEOUT | Таймаут чтения заголовка запроса | 2s                                                 |
| X509_CERT_PATH      | Путь до сертификата x509         | server.crt                                         |
| TLS_KEY_PATH        | Путь до ключа TLS                | server.key                                         |

## Клиент

### Эксплуатация

#### Конфигурация

Переменные необходимые для запуска приложения, указываются в файле `config.json`, файл должен быть расположен в той же директории, что и исполняемый:

| Параметр          | Описание                | По умолчанию |
|-------------------|-------------------------|--------------|
| db_file_path      | Путь до базы данных     | user.db      |
| db_client_timeout | Таймаут запроса клиента | 30s          |
| server_url        | URL сервера             |              |

#### Список доступных команд

* `gophkeeper register [username] [password]` - регистрация нового пользователя по логину и паролю;
* `gophkeeper login [username] [password]` - авторизация пользователя по логину и паролю;
* `gophkeeper create text [content]` - создать новые текстовые данные;
* `gophkeeper create binary [path-to-file]` - создать новые бинарные данные из файла;
* `gophkeeper create credentials [name] [login] [password]` - создать новый логин и пароль;
* `gophkeeper create bank-card [number] [valid-thru] [cvv] [(optional) card-holder]` - создать новую банковскую карту;
* `gophkeeper update text [id] [content]` - обновить существующие текстовые данные;
* `gophkeeper update binary [id] [path-to-file]` - обновить существующие бинарные данные;
* `gophkeeper update credentials --name --login --password [id]` - обновить существующие логин и пароль;
* `gophkeeper update bank-card --number --valid-thru --cvv --card-holder [id]` - обновить существующую банковскую карту;
* `gophkeeper show texts` - показать локальные текстовые данные;
* `gophkeeper show binaries` - показать локальные бинарные данные;
* `gophkeeper show credentials` - показать локальные логины и пароли;
* `gophkeeper show bank-cards` - показать локальные банковские карты;
* `gophkeeper sync texts` - синхронизировать (перезаписать) локальные текстовые данные;
* `gophkeeper sync binaries` - синхронизировать (перезаписать) локальные бинарные данные;
* `gophkeeper sync credentials` - синхронизировать (перезаписать) локальные логины и пароли;
* `gophkeeper sync bank-cards` - синхронизировать (перезаписать) локальные банковские карты;
* `gophkeeper sync all` - синхронизировать (перезаписать) все локальные данные;
* `gophkeeper help` - показать список всех команд или помощь для одной команды;

## Разработка

### Зависимости
* [Go 1.21](https://golang.org)
* [Make](https://www.gnu.org/software/make/)
* [swag](https://github.com/swaggo/swag)
* [godoc](https://cs.opensource.google/go/x/tools/+/master:godoc/)

### Команды

* `make all (default)` - последовательные запуск форматтеров, линтеров и тестов;
* `make build-client` - сборка бинарных файлов для cli приложения;
* `make clean` - очистка окружения;
* `make format` - форматирование исходного кода;
* `make generate-spec` - генерация swagger spec;
* `make godoc-get` - Получить документацию в формате html;
* `make godoc-run` - Запустить сервер документации;
* `make help` - вывод всех доступных команд;
* `make lint` - запуск статического анализа (линтеров);
* `make migration-down` - откат миграций бд сервера;
* `make migration-fix` - Фиксация миграций бд сервера (в случае неудачной установки/отката);
* `make migration-up` - Установка миграций бд сервера;
* `make setup` - установка виртуального окружения и всех зависимостей;
* `make test` - запуск юнит-тестов;

Список всех команд (Makefile targets) смотрите в [Makefile](Makefile).

### Изменение API

1. После изменения в сервере, добавьте [аннотации](https://github.com/swaggo/swag?tab=readme-ov-file#general-api-info);
2. Выполните команду `make generate-spec` для генерации `openapi`;

### [Документация](docs/arc42/)

Для всего проекта необходимо поддерживать документацию в формате [arc42](https://arc42.org/overview)
