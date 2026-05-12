# 🔧 Этап 2 — Use Case Architect

> Роль: System Analyst & Software Architect (Clean Architecture)  
> Входные данные: [User Stories + AC](../analytics/01-product-analyst.md)  
> Выходные данные: Use Cases (Interactors), Ports → [Этап 3](03-domain-entities.md)

---

## Маппинг: User Story → Use Case

| US ID | User Story | Use Case |
|---|---|---|
| US-01, US-02 | Создание комнаты + ссылка | `CreateRoom` |
| US-03 | Завершение сессии | `CloseRoom` |
| US-04 | Вход по ссылке | `JoinRoom` |
| US-05 | Видео/аудио (WebRTC) | `NegotiateWebRTC` |
| US-06 | Управление камерой | `ToggleCamera` |
| US-07 | Управление микрофоном | `ToggleMic` |
| US-08 | Отправка сообщений | `SendChatMessage` |
| US-09 | История чата | `GetChatHistory` |
| US-10 | Старт демонстрации экрана | `StartScreenShare` |
| US-11 | Стоп демонстрации | `StopScreenShare` |
| US-12 | Выход из комнаты | `LeaveRoom` |

---

## Карточки Use Cases

### UC-01: CreateRoom
**Актор:** Организатор | **Файл:** `internal/usecase/create_room.go`

**Request Model:**
```
organizerName: string
```

**Response Model (Success):**
```
roomID: string
inviteURL: string  — "/room/{roomID}"
organizerToken: string  — JWT
```

**Порты:**
- `RoomRepository.Save(room)`
- `TokenGenerator.NewToken(roomID, participantID, ORGANIZER)`

**Happy Path:**
1. Проверить organizerName ≠ пустой.
2. Сгенерировать уникальный roomID.
3. Создать Entity `Room` (status=ACTIVE).
4. Создать Entity `Participant` (role=ORGANIZER).
5. Сохранить → `RoomRepository.Save`.
6. Сгенерировать organizerToken → `TokenGenerator`.
7. Вернуть Response Model.

**Альтернативные сценарии:**
- [E1] organizerName пустой → `ErrEmptyName`
- [E2] RoomRepository.Save упал → wrapped storage error

---

### UC-02: JoinRoom
**Актор:** Участник | **Файл:** `internal/usecase/join_room.go`

**Request Model:** `roomID, participantName`

**Response Model:** `participantID, participantToken, currentParticipants[]`

**Порты:** `RoomRepository.FindByID`, `RoomRepository.AddParticipant`, `TokenGenerator`

**Happy Path:**
1. Найти комнату → `RoomRepository.FindByID`.
2. Проверить: `room.CanJoin()` (ACTIVE + кол-во < 4).
3. Создать `Participant` (role=GUEST).
4. `RoomRepository.AddParticipant`.
5. Сгенерировать token → `TokenGenerator`.
6. Вернуть token + список текущих участников.

**Альтернативные сценарии:**
- [E1] Комната не найдена → `ErrRoomNotFound`
- [E2] Комната закрыта → `ErrRoomClosed`
- [E3] Комната заполнена → `ErrRoomFull`
- [E4] participantName пустой → `ErrEmptyName`

---

### UC-03: CloseRoom
**Актор:** Организатор | **Файл:** `internal/usecase/room_lifecycle.go`

**Request Model:** `roomID, organizerToken`

**Порты:** `RoomRepository.FindByID`, `RoomRepository.UpdateStatus`, `SignalingGateway.BroadcastRoomClosed`

**Happy Path:**
1. Верифицировать organizerToken, проверить role=ORGANIZER.
2. Найти комнату, вызвать `room.Close(now)`.
3. `RoomRepository.UpdateStatus(CLOSED)`.
4. `SignalingGateway.BroadcastRoomClosed`.

**Альтернативные сценарии:**
- [E1] Токен не организатора → `ErrForbidden`
- [E2] Комната уже закрыта → `ErrRoomAlreadyClosed`

---

### UC-04: LeaveRoom
**Актор:** Участник | **Файл:** `internal/usecase/room_lifecycle.go`

**Request Model:** `roomID, participantToken`

**Порты:** `RoomRepository.RemoveParticipant`, `SignalingGateway.BroadcastParticipantLeft`

---

### UC-05: NegotiateWebRTC
**Актор:** Участник | **Файл:** `internal/usecase/negotiate_webrtc.go`

> Сигнальный Use Case. Сервер — только relay, WebRTC-сессия живёт в браузере.

**Request Model:** `roomID, participantToken, sdpOffer, toParticipantID`

**Порты:** `SignalingGateway.ForwardSDP`

---

### UC-06 / UC-07: ToggleCamera / ToggleMic
**Актор:** Участник | **Файл:** `internal/usecase/room_lifecycle.go`

**Request Model:** `roomID, participantToken, enabled: bool`

**Порты:** `SignalingGateway.BroadcastMediaState`

---

### UC-08: SendChatMessage
**Актор:** Участник | **Файл:** `internal/usecase/send_chat_message.go`

**Request Model:** `roomID, participantToken, text`

**Response Model:** `messageID, sentAt`

**Порты:** `RoomRepository.FindByID`, `ChatRepository.Save`, `SignalingGateway.BroadcastChatMessage`

**Бизнес-правило:** Отправка возможна только в ACTIVE комнате.

**Альтернативные сценарии:**
- [E1] text пустой → `ErrEmptyMessage`
- [E2] Комната не ACTIVE → `ErrRoomClosed`

---

### UC-09: GetChatHistory
**Актор:** Участник | **Файл:** `internal/usecase/send_chat_message.go`

**Request Model:** `roomID, participantToken`

**Response Model:** `messages: []ChatMessage`

**Порты:** `ChatRepository.FindByRoom`

---

### UC-10 / UC-11: StartScreenShare / StopScreenShare
**Актор:** Участник | **Файл:** `internal/usecase/screen_share.go`

**Request Model:** `roomID, participantToken`

**Порты:** `RoomRepository.FindByID`, `RoomRepository.SetScreenSharer`, `SignalingGateway.BroadcastScreenShareState`

**Бизнес-правило:** Одновременно только один шарер → `room.HasScreenSharer()`.

---

## Сводная таблица Use Cases

| ID | Use Case | Актор | Ключевые входы | Success Output | Ошибки |
|---|---|---|---|---|---|
| UC-01 | CreateRoom | Организатор | organizerName | roomID, inviteURL, token | E1, E2 |
| UC-02 | JoinRoom | Участник | roomID, name | token, participants | E1–E4 |
| UC-03 | CloseRoom | Организатор | roomID, token | — | E1, E2 |
| UC-04 | LeaveRoom | Участник | roomID, token | — | E1 |
| UC-05 | NegotiateWebRTC | Участник | sdpOffer, toID | — | E1 |
| UC-06 | ToggleCamera | Участник | enabled, token | — | E1 |
| UC-07 | ToggleMic | Участник | enabled, token | — | E1 |
| UC-08 | SendChatMessage | Участник | text, token | messageID | E1, E2 |
| UC-09 | GetChatHistory | Участник | roomID, token | []messages | E1 |
| UC-10 | StartScreenShare | Участник | roomID, token | — | E1, E2 |
| UC-11 | StopScreenShare | Участник | roomID, token | — | E1 |

---

## Порты (Interfaces)

| Интерфейс | Методы | UC |
|---|---|---|
| `RoomRepository` | Save, FindByID, UpdateStatus, AddParticipant, RemoveParticipant, SetScreenSharer | UC-01–04, UC-10–11 |
| `ChatRepository` | Save, FindByRoom | UC-08–09 |
| `SignalingGateway` | BroadcastRoomClosed, BroadcastParticipantLeft, BroadcastMediaState, BroadcastChatMessage, BroadcastScreenShareState, ForwardSDP | UC-03–11 |
| `TokenGenerator` | NewToken, Verify | UC-01–11 |

> Все интерфейсы объявлены в `internal/usecase/port/`. Реализации — в `internal/infra/`.

---

→ Следующий этап: [03-domain-entities.md](03-domain-entities.md)
