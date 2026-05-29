# Green API Gateway

HTTP-прокси для методов [Green API](https://green-api.com) (WhatsApp).  
Принимает запросы с `idInstance` и `apiTokenInstance` в теле — позволяет использовать любой инстанс без перезапуска сервера.

## Требования

- Go 1.25+

## Быстрый старт

**1. Склонировать репозиторий**

```bash
git clone <repo-url>
cd test_task_green_api
```

**2. Создать `.env`**

```bash
echo "ENV=local" > .env
```

**3. Запустить**

```bash
make run
```

Сервер поднимается на `http://localhost:8080`.  
Открой браузер — на главной странице (`/`) будет UI для тестирования всех методов.

## Конфигурация

Параметры сервера задаются в `configs/local.json`:

```json
{
  "server": {
    "port": 8080,
    "read_timeout": "5s",
    "write_timeout": "5s",
    "max_header_bytes": 1000000
  },
  "handler": {
    "allowed_cors_origins": "*"
  }
}
```

Переменные окружения (`.env`):

| Переменная | Описание              | Обязательная |
|------------|-----------------------|:------------:|
| `ENV`      | `local` / `dev` / `prod` | ✓         |

## Сборка

```bash
make build       # → bin/app
./bin/app -config ./configs/local.json
```

## API

Все методы принимают `POST` с JSON-телом. Поля `idInstance` и `apiTokenInstance` обязательны для каждого запроса.

### GET /health

```
GET /health
```

```json
{ "status": "ok" }
```

---

### POST /v1/green-api/settings

Настройки инстанса.

```json
{
  "idInstance": "1234567890",
  "apiTokenInstance": "your-token"
}
```

---

### POST /v1/green-api/state

Состояние авторизации инстанса.

```json
{
  "idInstance": "1234567890",
  "apiTokenInstance": "your-token"
}
```

---

### POST /v1/green-api/send-message

Отправка текстового сообщения.

```json
{
  "idInstance": "1234567890",
  "apiTokenInstance": "your-token",
  "chatId": "79001234567@c.us",
  "message": "Hello!"
}
```

Ответ:

```json
{ "idMessage": "3EB0C7C7D09787C7C030" }
```

---

### POST /v1/green-api/send-file-by-url

Отправка файла по URL.

```json
{
  "idInstance": "1234567890",
  "apiTokenInstance": "your-token",
  "chatId": "79001234567@c.us",
  "urlFile": "https://example.com/image.png",
  "fileName": "image.png"
}
```

Ответ:

```json
{ "idMessage": "3EB0C7C7D09787C7C030" }
```

## Структура проекта

```
.
├── cmd/app/
│   ├── main.go          # точка входа
│   └── web/
│       └── index.html   # UI (встроен в бинарник через embed)
├── configs/
│   └── local.json
├── internal/
│   ├── config/          # конфигурация
│   ├── domain/          # доменные типы
│   ├── modules/
│   │   └── green_api/   # handler · service · runner · interface
│   ├── runner/          # регистрация роутов
│   ├── server/http/     # HTTP-сервер
│   └── shared/
│       ├── dto/         # маппинг models ↔ domain
│       ├── middleware/  # CORS, panic recovery, access log
│       ├── response/    # запись HTTP-ответов
│       ├── server_error/
│       └── transport_error/
├── models/              # сгенерировано из swagger.yaml (не редактировать)
├── pkg/
│   ├── client/green_api/ # HTTP-клиент Green API
│   └── logger/
├── swagger.yaml
└── Makefile
```

## Генерация моделей

Модели в `models/` генерируются из `swagger.yaml` с помощью [go-swagger](https://github.com/go-swagger/go-swagger):

```bash
make generate-models
```

При первом запуске бинарник swagger скачается автоматически в `bin/swagger`.
