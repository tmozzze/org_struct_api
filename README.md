
# OrgStruct API

## Запуск проекта

### Предварительные требования
- Установленный Docker и Docker Compose.
- Наличие файла `.env` (создайте его на основе `.env.example`).

### Быстрый старт
1. Соберите и запустите контейнеры:
    в терминале bash
   ```sh
   make up
   ```
   *Эта команда поднимет базу данных, применит миграции и запустит API сервер.*

2. Откройте Swagger-документацию:
   ```
   http://localhost:8080/swagger/index.html
   ```
3. Пользуйтесь :)

## Структура проекта

```text
.
├── cmd/
│   └── api/
│       └── main.go            # Точка входа в приложение
├── database/
│   └── migrations/            # SQL миграции
├── docs/                      # Swagger документация
├── internal/
│   ├── config/                # Загрузка конфигурации
│   ├── domain/                # Доменные модели
│   ├── handler/
│   │   └── http/              # HTTP слой
│   ├── repository/
│   │   └── postgres/          # Реализация репозиториев (GORM)
│   └── service/               # Бизнес-логика
├── pkg/
│   └── database/              # Хелперы для подключения к БД и миграций
├── docker-compose.yml         # Конфигурация Docker
├── Dockerfile                 
├── Makefile                   # Команды для автоматизации
├── .env.example               # Пример переменных окружения
└── go.mod                     # Зависимости проекта
```

## Технологический стек
- **Language:** Go
- **Database:** PostgreSQL
- **ORM:** GORM
- **Migrations:** Goose
- **Logging:** slog
- **Documentation:** Swagger (swaggo)
- **Containerization:** Docker & Docker Compose

## Команды управления (Makefile)

### Разработка и запуск
- `make run` — запуск приложения локально.
- `make test` — запуск всех тестов.
- `make swagger-gen` — генерация Swagger документации.
- `make lint` — запуск линтера (требуется golangci-lint).

### Docker
- `make up` — сборка и запуск проекта в Docker (в фоне).
- `make down` — остановка контейнеров.
- `make down-and-clean` — остановка и удаление всех данных (volumes).

### Миграции (Goose)
- `make migrate-up` — применить все миграции.
- `make migrate-down` — откатить последнюю миграцию.
- `make migrate-status` — проверить статус миграций.
- `make migrate-reset` — откатить все миграции.
