# ⚙️ Этап 4 — System Design Architect

> Роль: Tech Lead  
> Входные данные: [Domain Entities + Ports](03-domain-entities.md)  
> Выходные данные: TDR, стек, схема данных, структура пакетов → [Этап 5](05-devops-workflow.md)

---

## Внешние границы системы

| Компонент | Выбор | Обоснование |
|---|---|---|
| **HTTP API** | REST (chi) | CRUD-операции: создать/войти/выйти/история чата |
| **Real-time** | WebSocket (gorilla/websocket) | Сигнальный канал: SDP relay, chat broadcast, события комнаты |
| **База данных** | Redis 7 | Состояние эфемерное (сессия = комната). In-memory, TTL 24 ч, автоочистка. PostgreSQL избыточен |
| **Брокер** | нет | 3–4 участника, одна нода — WebSocket fan-out через goroutine hub достаточен |
| **WebRTC** | Mesh P2P | SFU/MCU не нужны для ≤4 участников. Сервер — только signaling relay |
| **STUN/TURN** | Coturn (внешний) | Обход NAT для P2P WebRTC. Конфигурируется через env |
| **Авторизация** | JWT (HMAC-SHA256) | Stateless, не требует хранилища. Payload: roomID + participantID + role |

---

## Ports → Adapters

| Port (интерфейс) | Adapter (реализация) | Файл |
|---|---|---|
| `RoomRepository` | `redis.RoomRepository` | `internal/infra/redis/room_repository.go` |
| `ChatRepository` | `redis.ChatRepository` | `internal/infra/redis/chat_repository.go` |
| `SignalingGateway` | `ws.Hub` (goroutine fan-out) | `internal/infra/ws/hub.go` |
| `TokenGenerator` | `token.JWTGenerator` | `internal/infra/token/jwt_generator.go` |

---

## Redis-схема (ephemeral state)

```
room:{roomID}        STRING  JSON(Room)            TTL: 24h
chat:{roomID}:messages  LIST  [ JSON(ChatMessage) ]  TTL: 24h
```

> Всё хранится в одном JSON-блоке на комнату. Индексирование не нужно — доступ только по roomID.

---

## HTTP API

| Метод | Путь | Use Case | Auth |
|---|---|---|---|
| `POST` | `/api/v1/rooms` | CreateRoom | нет |
| `POST` | `/api/v1/rooms/:id/join` | JoinRoom | нет |
| `DELETE` | `/api/v1/rooms/:id` | CloseRoom | JWT (ORGANIZER) |
| `DELETE` | `/api/v1/rooms/:id/leave` | LeaveRoom | JWT |
| `GET` | `/api/v1/rooms/:id/chat` | GetChatHistory | JWT |
| `POST` | `/api/v1/rooms/:id/chat` | SendChatMessage | JWT |
| `GET` | `/ws/:roomID?token=...` | WebSocket signaling | JWT (query param) |

---

## WebSocket сигнальный протокол

Все сообщения — JSON. Тип определяется полем `"type"`.

**Клиент → Сервер:**

| type | Поля | Описание |
|---|---|---|
| `sdp` | `toId, sdp` | Переслать SDP offer/answer конкретному участнику |
| `ice` | `toId, candidate` | Переслать ICE candidate |

**Сервер → Клиент (broadcast):**

| type | Поля | Описание |
|---|---|---|
| `room_closed` | `roomId` | Комната завершена организатором |
| `participant_left` | `participantId` | Участник покинул комнату |
| `media_state` | `participantId, camera, mic` | Изменение состояния камеры/микрофона |
| `chat_message` | `id, senderName, text, sentAt` | Новое сообщение в чате |
| `screen_share_state` | `participantId, active` | Старт/стоп демонстрации экрана |

**Сервер → Клиент (unicast):**

| type | Поля | Описание |
|---|---|---|
| `sdp` | `sdp` | Входящий SDP от другого участника |
| `ice` | `candidate` | Входящий ICE candidate |

---

## Структура пакетов

```
meet/
├── cmd/
│   └── server/
│       └── main.go                      # точка входа, ручной DI
├── internal/
│   ├── domain/
│   │   ├── room.go                      # Entity Room + RoomStatus
│   │   ├── participant.go               # Entity Participant + ParticipantRole
│   │   ├── chat_message.go              # Entity ChatMessage
│   │   └── errors.go                    # все ErrXxx
│   ├── usecase/
│   │   ├── port/
│   │   │   ├── room_repository.go       # Port: RoomRepository
│   │   │   ├── chat_repository.go       # Port: ChatRepository
│   │   │   ├── signaling_gateway.go     # Port: SignalingGateway
│   │   │   └── token_generator.go       # Port: TokenGenerator
│   │   ├── create_room.go
│   │   ├── join_room.go
│   │   ├── room_lifecycle.go            # CloseRoom, LeaveRoom, ToggleCamera, ToggleMic
│   │   ├── send_chat_message.go         # SendChatMessage + GetChatHistory
│   │   └── screen_share.go             # StartScreenShare, StopScreenShare
│   ├── adapter/
│   │   ├── http/
│   │   │   ├── handler/
│   │   │   │   ├── room_handler.go
│   │   │   │   └── chat_handler.go
│   │   │   ├── dto/
│   │   │   │   └── room_dto.go
│   │   │   └── router.go
│   │   └── ws/
│   │       └── ws_handler.go
│   └── infra/
│       ├── redis/
│       │   ├── room_repository.go
│       │   └── chat_repository.go
│       ├── ws/
│       │   └── hub.go
│       └── token/
│           └── jwt_generator.go
├── config/
│   └── config.go
├── docs/
├── .github/workflows/
│   └── ci.yml
├── .golangci.yml
├── Makefile
├── Dockerfile
└── docker-compose.yml
```

---

## Библиотеки

| Категория | Библиотека | Версия | Обоснование |
|---|---|---|---|
| HTTP-роутер | `go-chi/chi/v5` | v5.0.x | Лёгкий, stdlib-совместимый |
| WebSocket | `gorilla/websocket` | v1.5.x | Стандарт де-факто |
| Redis | `redis/go-redis/v9` | v9.x | Актуальная версия, ctx-first |
| JWT | `golang-jwt/jwt/v5` | v5.x | Минималистичный, активно поддерживается |
| UUID | `google/uuid` | v1.x | Стандарт |
| Конфигурация | `kelseyhightower/envconfig` | v1.x | 12-factor, без лишнего |
| Логирование | `log/slog` (stdlib) | Go 1.21+ | Ничего лишнего для небольшого сервиса |
| Тесты | `stretchr/testify` | v1.9.x | assertions |
| Моки | `uber-go/mock` | v0.4.x | Генерация mock-портов |

---

## TDR: Meet — Signaling Server

**Версия:** 1.0  
**Стек:** Go 1.22 · chi · gorilla/websocket · Redis 7 · JWT · Docker

### Архитектурные решения

| # | Решение | Альтернативы | Причина выбора |
|---|---|---|---|
| 1 | Clean Architecture (domain → usecase → adapter → infra) | MVC, Hexagonal | Изоляция бизнес-логики, тестируемость через моки портов |
| 2 | WebRTC Mesh P2P, сервер = signaling relay | SFU (mediasoup), MCU | 3–4 участника — SFU избыточен. Mesh дешевле в инфре |
| 3 | Redis (ephemeral) вместо PostgreSQL | PostgreSQL, SQLite | Комнаты живут 24 ч. Redis TTL = встроенная очистка, нет миграций |
| 4 | In-memory WebSocket Hub | Redis Pub/Sub, NATS | Одна нода, нет смысла в брокере. Горутина-хаб = 0 внешних зависимостей |
| 5 | JWT stateless | Redis sessions | Нет хранилища сессий, горизонтальное масштабирование при необходимости |
| 6 | Один бинарник: HTTP + WS | Отдельные сервисы | Для масштаба 3–4 участника это оверинжиниринг |

### Известные ограничения
- Одна нода: при росте нагрузки Hub нужно заменить на Redis Pub/Sub
- Нет записи звонка (Won't have в v1)
- TURN-сервер — внешняя зависимость, не входит в этот репозиторий

---

→ Следующий этап: [05-devops-workflow.md](05-devops-workflow.md)
