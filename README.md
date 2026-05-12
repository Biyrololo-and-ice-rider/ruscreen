# 🎥 Meet

Минималистичный сервис видеоконференций для 3–4 человек.  
WebRTC mesh P2P + сигнальный сервер на Go + текстовый чат + демонстрация экрана.

## Стек

| Компонент | Технология |
|---|---|
| Язык | Go 1.22 |
| HTTP / WebSocket | chi + gorilla/websocket |
| Хранилище | Redis 7 (ephemeral state) |
| Авторизация | JWT (HMAC-SHA256) |
| WebRTC | Mesh P2P, сервер — только signaling relay |
| Контейнеризация | Docker + docker-compose |

## Быстрый старт

```bash
# 1. Поднять Redis
make docker-up

# 2. Запустить сервер
make run

# 3. Открыть в браузере
open http://localhost:8080
```

## Команды

```bash
make build            # собрать бинарник → bin/meet
make run              # запустить сервер
make dev              # docker-compose up + go run
make test             # unit + integration тесты
make test-unit        # только unit (domain + usecase)
make test-integration # только integration (infra + adapter)
make lint             # golangci-lint
make docker-up        # поднять Redis в Docker
make docker-down      # остановить контейнеры
```

## Структура проекта

```
meet/
├── cmd/server/           # точка входа
├── internal/
│   ├── domain/           # Entities + Enterprise Business Rules
│   ├── usecase/          # Interactors + порты (интерфейсы)
│   │   └── port/
│   ├── adapter/          # HTTP handlers + WebSocket handler
│   │   ├── http/
│   │   └── ws/
│   └── infra/            # реализации портов (Redis, JWT, WS Hub)
│       ├── redis/
│       ├── token/
│       └── ws/
├── config/
├── migrations/           # (зарезервировано, сейчас Redis TTL)
├── docs/
│   ├── analytics/        # Product Analyst (User Stories, бэклог)
│   └── architecture/     # Use Cases, Entities, System Design, DevOps, Sprint
└── .github/workflows/
```

## Документация

| Документ | Описание |
|---|---|
| [docs/analytics/01-product-analyst.md](docs/analytics/01-product-analyst.md) | User Stories, AC, бэклог |
| [docs/architecture/02-use-cases.md](docs/architecture/02-use-cases.md) | Use Cases (Clean Architecture) |
| [docs/architecture/03-domain-entities.md](docs/architecture/03-domain-entities.md) | Domain Layer |
| [docs/architecture/04-system-design.md](docs/architecture/04-system-design.md) | Tech Design Record |
| [docs/architecture/05-devops-workflow.md](docs/architecture/05-devops-workflow.md) | Git Flow, CI, линтер |
| [docs/architecture/06-sprint-backlog.md](docs/architecture/06-sprint-backlog.md) | Sprint 1 Backlog |

## API

| Метод | Путь | Описание |
|---|---|---|
| `POST` | `/api/v1/rooms` | Создать комнату |
| `POST` | `/api/v1/rooms/:id/join` | Войти в комнату |
| `DELETE` | `/api/v1/rooms/:id` | Закрыть комнату (организатор) |
| `GET` | `/api/v1/rooms/:id/chat` | История чата |
| `WS` | `/ws/:roomID` | WebSocket signaling |

## Архитектура

Проект следует **Clean Architecture** (Uncle Bob):

```
domain → usecase → adapter → infra
```

- `domain/` не импортирует ничего кроме stdlib
- `usecase/` работает только через интерфейсы из `usecase/port/`
- `infra/` — единственное место, где живут Redis, JWT, WebSocket Hub
- `adapter/` переводит HTTP/WS ↔ Request/Response Models

## Contributing

См. [CONTRIBUTING.md](CONTRIBUTING.md)
