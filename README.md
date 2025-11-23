# Avito PR Test Service

Сервис для управления Pull Request'ами с автоматическим назначением ревьюеров из команды.

## Структура проекта

### Сервисы

**app** - основное HTTP приложение
- **Порт**: 8080 (настраивается через `HTTP_PORT`)
- Зависит от PostgreSQL и миграций
- Авторестарт при ошибках

**pg** - база данных PostgreSQL
- **Порт**: 5432 (настраивается через `PG_PORT`)
- Volume для данных (`postgres_volume`)
- Health check

**migrator** - сервис для выполнения миграций базы данных
- Запускается автоматически перед стартом приложения
- Использует goose для миграций

### Dockerfile

- **config/Dockerfile** - финальный образ основного сервиса
- **config/migration.Dockerfile** - образ для миграций

## API Эндпоинты

### Pull Requests

- `POST /pullRequest/create` - создание нового PR с автоматическим назначением ревьюеров
- `POST /pullRequest/merge` - слияние PR
- `POST /pullRequest/reassign` - переназначение ревьюера

### Teams

- `POST /team/add` - создание команды
- `GET /team/get` - получение информации о команде

### Users

- `POST /users/setIsActive` - установка статуса активности пользователя
- `GET /users/getReview` - получение списка PR, назначенных на ревью пользователю

### Statistics

- `GET /stats` - получение статистики по системе

Пример ответа:
```json
{
    "total_teams": 2,
    "total_users": 5,
    "active_users": 4,
    "total_prs": 7,
    "open_prs": 7,
    "merged_prs": 0,
    "total_assignments": 9
}
```

## Конфигурация

### Переменные окружения

Приложение использует переменные окружения для конфигурации. Основные переменные:

```bash
# HTTP сервер
HTTP_HOST=0.0.0.0          # Хост HTTP сервера
HTTP_PORT=8080            # Порт HTTP сервера

# PostgreSQL
PG_DSN=postgres://postgres:postgres@pg:5432/avito_pr?sslmode=disable  # DSN для подключения к БД

# Для docker-compose (опционально)
PG_USER=postgres           # Пользователь PostgreSQL
PG_PASSWORD=postgres       # Пароль PostgreSQL
PG_DATABASE_NAME=avito_pr # Имя базы данных
PG_PORT=5432               # Порт PostgreSQL
MIGRATION_DIR=migrations   # Директория с миграциями
```

### Логирование

Логирование настроено через `go.uber.org/zap`:
- Логи выводятся в stdout и в файл `logs/app.log`
- Уровень логирования настраивается через флаг `-level` (по умолчанию `info`)
- Доступные уровни: `debug`, `info`, `warn`, `error`

## Настройка окружения

1. Скопируйте пример переменных окружения (если есть `.env.example`):
   ```bash
   cp .env.example .env
   ```

2. Заполните переменные окружения в файле `.env` при необходимости

3. Убедитесь, что директория `logs/` существует:
   ```bash
   mkdir -p logs
   ```

## Запуск сервиса

### С помощью Docker Compose (рекомендуется)

```bash
# Сборка образов
make docker-build

# Запуск всех сервисов
make docker-up

# Остановка и очистка
make docker-clean
```

### Локальный запуск

1. Убедитесь, что PostgreSQL запущен и доступен

2. Выполните миграции:
   ```bash
   # Используя goose из bin/
   ./bin/goose -dir migrations postgres "your-dsn" up
   ```

3. Установите переменные окружения:
   ```bash
   export HTTP_HOST=0.0.0.0
   export HTTP_PORT=8080
   export PG_DSN="postgres://user:password@localhost:5432/avito_pr?sslmode=disable"
   ```

4. Запустите приложение:
   ```bash
   make build
   ./bin/app
   ```

   Или напрямую:
   ```bash
   go run ./cmd/http-server/main.go
   ```

## Разработка

### Makefile команды

- `docker-up` - Запускает Docker-контейнеры (docker compose up -d)
- `docker-clean` - Останавливает и удаляет Docker-контейнеры и связанные с ними образы
- `docker-build` - Собирает Docker-образы для всех сервисов
- `test` - Запускает тесты Go с флагами для обнаружения гонок и сбора покрытия кода
- `lint` - Запускает статический анализатор кода golangci-lint для всего проекта
- `deps` - Управляет зависимостями Go-модулей: загружает, проверяет и очищает их
- `build` - Компилирует приложение и сохраняет исполняемый файл в `bin/app`
- `clean` - Удаляет артефакты сборки (директория `bin/` и файл `cover.out`)
- `mock` - Генерирует моки для интерфейсов, используя mockery (конфигурация в `config/.mockery.yaml`)

### Тестирование

```bash
# Запуск всех тестов
make test

# Запуск тестов с покрытием
go test -coverprofile=cover.out ./...
go tool cover -html=cover.out
```

### Линтинг

```bash
# Запуск линтера
make lint

# Или напрямую
golangci-lint run ./...
```

Конфигурация линтера находится в `.golangci.yml`.

### Генерация моков

```bash
make mock
```

Конфигурация mockery находится в `config/.mockery.yaml`.

## Структура проекта

```
.
├── cmd/
│   └── http-server/      # Точка входа приложения
├── internal/
│   ├── app/              # Инициализация приложения
│   ├── config/           # Конфигурация (HTTP, PostgreSQL)
│   ├── http-server/      # HTTP handlers и middleware
│   ├── logger/           # Логирование
│   ├── model/            # Модели данных
│   ├── repository/       # Слой доступа к данным
│   └── service/          # Бизнес-логика
├── config/               # Dockerfile и конфигурация миграций
├── migrations/          # SQL миграции (goose)
├── .golangci.yml        # Конфигурация golangci-lint
├── docker-compose.yml   # Docker Compose конфигурация
└── Makefile             # Команды для разработки
```

## Дополнительные задания

- ✅ Реализован эндпоинт статистики `GET /stats`
- ✅ Настроен линтер golangci-lint
- ✅ Настроен CI/CD через GitHub Actions (`.github/workflows/go.yaml`)

## Технологии

- **Go 1.24.1** - язык программирования
- **PostgreSQL** - база данных
- **pgx/v5** - драйвер PostgreSQL
- **goose** - миграции базы данных
- **zap** - структурированное логирование
- **golangci-lint** - статический анализ кода
- **mockery** - генерация моков для тестирования
- **testify** - библиотека для тестирования
