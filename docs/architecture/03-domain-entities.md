# 🏗 Этап 3 — Domain Entities

> Роль: Software Architect (Clean Architecture by Uncle Bob)  
> Входные данные: [Use Cases + Ports](02-use-cases.md)  
> Выходные данные: Domain Layer → [Этап 4](04-system-design.md)

> ⚠️ Это **Clean Architecture**, не DDD. Entities здесь — структуры + методы валидации Enterprise Business Rules, независимые от БД и фреймворков.

---

## Entities

| Entity | Описание | Файл |
|---|---|---|
| `Room` | Конференц-комната с участниками и состоянием | `internal/domain/room.go` |
| `Participant` | Участник комнаты | `internal/domain/participant.go` |
| `ChatMessage` | Сообщение в чате | `internal/domain/chat_message.go` |

---

## Entity: Room

**Поля:**

| Поле | Тип | Описание |
|---|---|---|
| ID | string | Уникальный короткий ID (nanoid/uuid[:8]) |
| Status | RoomStatus | `ACTIVE` / `CLOSED` |
| Participants | []Participant | Встроены — не имеют смысла без Room |
| ScreenSharerID | *string | ID участника, шарящего экран (nil = никто) |
| CreatedAt | time.Time | Момент создания |
| ClosedAt | *time.Time | Момент закрытия (nil если ACTIVE) |

**Конструктор:** `NewRoom(id string, createdAt time.Time) (Room, error)`

**Enterprise Business Rules:**

| Метод | Инвариант |
|---|---|
| `CanJoin() bool` | Status == ACTIVE && len(Participants) < 4 |
| `IsActive() bool` | Status == ACTIVE |
| `HasScreenSharer() bool` | ScreenSharerID != nil |
| `Close(closedAt time.Time) error` | Только если ACTIVE, иначе `ErrRoomAlreadyClosed` |

---

## Entity: Participant

**Поля:**

| Поле | Тип | Описание |
|---|---|---|
| ID | string | UUID |
| Name | string | Отображаемое имя (непустое, ≤50 символов) |
| Role | ParticipantRole | `ORGANIZER` / `GUEST` |
| CameraEnabled | bool | Состояние камеры (default: true) |
| MicEnabled | bool | Состояние микрофона (default: true) |
| JoinedAt | time.Time | Момент входа |

**Конструктор:** `NewParticipant(id, name string, role ParticipantRole, joinedAt time.Time) (Participant, error)`

**Enterprise Business Rules:**

| Метод | Инвариант |
|---|---|
| `IsOrganizer() bool` | Role == ORGANIZER |
| `validateName() error` | Name непустой, len([]rune) ≤ 50 |

---

## Entity: ChatMessage

**Поля:**

| Поле | Тип | Описание |
|---|---|---|
| ID | string | UUID |
| RoomID | string | Ссылка через ID |
| SenderID | string | ID участника-отправителя |
| SenderName | string | Денормализовано для истории |
| Text | string | Текст (непустой, ≤2000 символов) |
| SentAt | time.Time | Время отправки |

**Конструктор:** `NewChatMessage(id, roomID, senderID, senderName, text string, sentAt time.Time) (ChatMessage, error)`

**Enterprise Business Rules:** `validateText() error` — Text непустой, len([]rune) ≤ 2000.

---

## Enterprise Business Rules — сводная таблица

| Правило | Entity | Ошибка |
|---|---|---|
| Комната не принимает участника при Status=CLOSED | Room | `ErrRoomClosed` |
| Комната не принимает >4 участников | Room | `ErrRoomFull` |
| Нельзя закрыть уже закрытую комнату | Room | `ErrRoomAlreadyClosed` |
| Имя участника не пустое | Participant | `ErrEmptyName` |
| Имя участника ≤50 символов | Participant | `ErrNameTooLong` |
| Текст сообщения не пустой | ChatMessage | `ErrEmptyMessage` |
| Текст сообщения ≤2000 символов | ChatMessage | `ErrMessageTooLong` |

---

## Связи между Entities

| От | К | Тип | Реализация | Комментарий |
|---|---|---|---|---|
| Room | Participant | 1:N | `Room.Participants []Participant` | Встроен: участник без комнаты не существует |
| ChatMessage | Room | N:1 | `ChatMessage.RoomID string` | Только ID — нет прямой ссылки |
| ChatMessage | Participant | N:1 | `ChatMessage.SenderID string` | Только ID; SenderName денормализован |

---

## Все ошибки Domain Layer

```go
// internal/domain/errors.go

var (
    ErrRoomNotFound      = errors.New("room not found")
    ErrRoomClosed        = errors.New("room is closed")
    ErrRoomFull          = errors.New("room is full (max 4 participants)")
    ErrRoomAlreadyClosed = errors.New("room is already closed")

    ErrEmptyName   = errors.New("name must not be empty")
    ErrNameTooLong = errors.New("name must not exceed 50 characters")

    ErrEmptyMessage   = errors.New("message must not be empty")
    ErrMessageTooLong = errors.New("message must not exceed 2000 characters")

    ErrForbidden         = errors.New("forbidden: insufficient role")
    ErrInvalidToken      = errors.New("invalid or expired token")
    ErrScreenShareActive = errors.New("screen share already active by another participant")
)
```

---

→ Следующий этап: [04-system-design.md](04-system-design.md)
