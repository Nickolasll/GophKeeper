# GophKeeper
Менеджер паролей GophKeeper

## Сервер

### [API спецификация](docs/api/swagger.json)
У запущенного сервиса есть endpoint с актуальной спецификацией на `/api/v1/swagger.json`. Его можно использовать в Swagger UI или Postman для импорта.

### Эксплуатация

#### Логи
Пишутся в формате json. Пример лога http-запроса:

```json
{"content length":0,"duration_ms":84,"level":"info","method":"POST","msg":"Request info","status":200,"time":"2024-03-05T17:50:57+03:00","uri":"/api/v1/auth/register"}
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

TODO

## Разработка

### Зависимости
* [Go 1.21](https://golang.org)
* [Make](https://www.gnu.org/software/make/)
* [swag](https://github.com/swaggo/swag)
* [godoc](https://cs.opensource.google/go/x/tools/+/master:godoc/)

### Команды

* `make setup` - установка виртуального окружения и всех зависимостей;
* `make format` - форматирование исходного кода;
* `make lint` - запуск статического анализа (линтеров);
* `make test` - запуск юнит-тестов;
* `make client-build` - сборка исполняемых файлов клиента;
* `make generate-spec` - кодогенерация по openapi-спецификации;
* `make help` - вывод всех доступных команд;

Список всех команд (Makefile targets) смотрите в [Makefile](Makefile).

### Изменение API

1. После изменения в сервере, добавьте [аннотации](https://github.com/swaggo/swag?tab=readme-ov-file#general-api-info);
2. Выполните команду `make generate-spec` для генерации `openapi`;

### [Документация](docs/arc42/)

Для всего проекта необходимо поддерживать документацию в формате [arc42](https://arc42.org/overview)
