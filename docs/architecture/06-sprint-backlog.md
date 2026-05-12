# 🗂 Этап 6 — Sprint Planner

> Роль: Team Lead + Scrum Master  
> Входные данные: все предыдущие этапы  
> Выходные данные: Sprint 1 Backlog, Story Points, CSV для Jira/Linear

---

## Ёмкость команды

| Разработчик | Роль | SP / спринт |
|---|---|---|
| Dev 1 | Senior Go (WebSocket, infra) | 20 SP |
| Dev 2 | Go Backend | 18 SP |
| Dev 3 | Frontend (WebRTC, UI) | 16 SP |
| **Итого** | | **54 SP** |

---

## Sprint Backlog

### Слой: Инфраструктура проекта

| ID | Задача | SP | Роль | Приоритет | Зависит от |
|---|---|---|---|---|---|
| T-01 | Init репо: go.mod, структура папок, Dockerfile, docker-compose (Redis) | 2 | Senior | 🔴 Must | — |
| T-02 | Настройка golangci-lint + .golangci.yml | 1 | Senior | 🔴 Must | T-01 |
| T-03 | GitHub Actions CI (lint + test-unit + test-integration + build + Redis service) | 3 | Senior | 🔴 Must | T-02 |

### Слой: Domain

| ID | Задача | SP | Роль | Приоритет | Зависит от |
|---|---|---|---|---|---|
| T-04 | Entity Room: поля, NewRoom, CanJoin, IsActive, HasScreenSharer, Close + инварианты | 3 | Dev | 🔴 Must | T-01 |
| T-05 | Entity Participant: поля, NewParticipant, IsOrganizer, validateName | 2 | Dev | 🔴 Must | T-01 |
| T-06 | Entity ChatMessage: поля, NewChatMessage, validateText | 2 | Dev | 🔴 Must | T-01 |
| T-07 | domain/errors.go — все ErrXxx | 1 | Dev | 🔴 Must | T-04, T-05, T-06 |
| T-08 | Unit-тесты: Room, Participant, ChatMessage (покрытие ≥ 80%) | 3 | Dev | 🔴 Must | T-07 |

### Слой: Ports (интерфейсы)

| ID | Задача | SP | Роль | Приоритет | Зависит от |
|---|---|---|---|---|---|
| T-09 | Порты: RoomRepository, ChatRepository | 1 | Senior | 🔴 Must | T-07 |
| T-10 | Порты: SignalingGateway, TokenGenerator | 1 | Senior | 🔴 Must | T-07 |

### Слой: Infra (реализации портов)

| ID | Задача | SP | Роль | Приоритет | Зависит от |
|---|---|---|---|---|---|
| T-11 | JWT TokenGenerator: NewToken, Verify | 3 | Dev | 🔴 Must | T-10 |
| T-12 | Redis RoomRepository: Save, FindByID, UpdateStatus, AddParticipant, RemoveParticipant, SetScreenSharer | 5 | Dev | 🔴 Must | T-09 |
| T-13 | Redis ChatRepository: Save, FindByRoom | 3 | Dev | 🔴 Must | T-09 |
| T-14 | Integration-тесты: RoomRepository + ChatRepository (тег integration) | 3 | Dev | 🔴 Must | T-12, T-13 |
| T-15 | WebSocket Hub (SignalingGateway): Register/Unregister, broadcast fan-out, ForwardSDP | 5 | Senior | 🔴 Must | T-10 |

### Слой: Use Cases

| ID | Задача | SP | Роль | Приоритет | Зависит от |
|---|---|---|---|---|---|
| T-16 | UseCase CreateRoom + unit-тесты (mock порты) | 3 | Dev | 🔴 Must | T-09, T-10 |
| T-17 | UseCase JoinRoom + unit-тесты | 3 | Dev | 🔴 Must | T-09, T-10 |
| T-18 | UseCase CloseRoom + LeaveRoom + unit-тесты | 3 | Dev | 🔴 Must | T-09, T-10 |
| T-19 | UseCase ToggleCamera + ToggleMic + unit-тесты | 2 | Dev | 🔴 Must | T-10 |
| T-20 | UseCase NegotiateWebRTC (SDP relay) + unit-тесты | 3 | Senior | 🔴 Must | T-10 |
| T-21 | UseCase SendChatMessage + GetChatHistory + unit-тесты | 3 | Dev | 🟡 Should | T-09, T-10 |
| T-22 | UseCase StartScreenShare + StopScreenShare + unit-тесты | 3 | Dev | 🟡 Should | T-09, T-10 |

### Слой: Adapters (HTTP + WebSocket)

| ID | Задача | SP | Роль | Приоритет | Зависит от |
|---|---|---|---|---|---|
| T-23 | HTTP Handler: POST /rooms (CreateRoom) + DTO + JWT middleware | 3 | Dev | 🔴 Must | T-16, T-11 |
| T-24 | HTTP Handler: POST /rooms/:id/join + GET /rooms/:id/chat | 3 | Dev | 🔴 Must | T-17, T-21 |
| T-25 | WS Handler: upgrade + роутинг SDP/ICE сигналов | 5 | Senior | 🔴 Must | T-15, T-20 |
| T-26 | HTTP Handler: DELETE /rooms/:id (CloseRoom) + /leave | 2 | Dev | 🟡 Should | T-18, T-23 |

### Слой: Frontend

| ID | Задача | SP | Роль | Приоритет | Зависит от |
|---|---|---|---|---|---|
| T-27 | Фронтенд: базовый WebRTC клиент — создать/войти, видео/аудио mesh | 8 | Frontend | 🔴 Must | T-25 |
| T-28 | Фронтенд: чат UI — отправка, получение в реальном времени | 3 | Frontend | 🟡 Should | T-27 |
| T-29 | Фронтенд: демонстрация экрана — start/stop UI | 3 | Frontend | 🟡 Should | T-27 |

---

## Итог Sprint 1

| Категория | Задач | SP |
|---|---|---|
| 🔴 Must have | 22 | 46 |
| 🟡 Should have | 7 | 17 |
| **Запланировано в спринт** | **22** (Must) | **46** |
| Запас команды | | ~8 SP |

> Should have задачи (T-21, T-22, T-26, T-28, T-29) подтягиваются по ходу спринта по мере закрытия Must задач.

---

## Граф зависимостей

```
T-01 (init)
  ├── T-02 (lint)
  │     └── T-03 (CI)
  ├── T-04, T-05, T-06 (entities)
  │     └── T-07 (errors)
  │           ├── T-08 (entity tests)
  │           ├── T-09 (ports repo)
  │           └── T-10 (ports signals)
  │                 ├── T-11 (JWT)
  │                 ├── T-12 (Redis room) ──┐
  │                 ├── T-13 (Redis chat) ──┴── T-14 (integration tests)
  │                 ├── T-15 (WS Hub)
  │                 ├── T-16 (UC CreateRoom) ──── T-23 (HTTP Create)
  │                 ├── T-17 (UC JoinRoom)  ──── T-24 (HTTP Join)
  │                 ├── T-18 (UC Close/Leave) ── T-26 (HTTP Close)
  │                 ├── T-19 (UC ToggleMedia)
  │                 └── T-20 (UC WebRTC) ──────── T-25 (WS Handler)
  │                                                    └── T-27 (Frontend)
  │                                                          ├── T-28 (Chat UI)
  │                                                          └── T-29 (Screen UI)
```

---

## CSV для импорта в Jira / Linear / Trello

```csv
ID,Summary,Layer,StoryPoints,Role,Priority,DependsOn
T-01,Init repo structure Dockerfile docker-compose,infra,2,Senior,Must,
T-02,golangci-lint setup,ci,1,Senior,Must,T-01
T-03,GitHub Actions CI pipeline,ci,3,Senior,Must,T-02
T-04,Entity Room,domain,3,Dev,Must,T-01
T-05,Entity Participant,domain,2,Dev,Must,T-01
T-06,Entity ChatMessage,domain,2,Dev,Must,T-01
T-07,domain errors ErrXxx,domain,1,Dev,Must,T-04
T-08,Unit tests all entities,domain,3,Dev,Must,T-07
T-09,Ports RoomRepository ChatRepository,usecase,1,Senior,Must,T-07
T-10,Ports SignalingGateway TokenGenerator,usecase,1,Senior,Must,T-07
T-11,Infra JWT TokenGenerator,infra,3,Dev,Must,T-10
T-12,Infra Redis RoomRepository,infra,5,Dev,Must,T-09
T-13,Infra Redis ChatRepository,infra,3,Dev,Must,T-09
T-14,Integration tests Redis repos,infra,3,Dev,Must,T-12
T-15,Infra WebSocket Hub SignalingGateway,infra,5,Senior,Must,T-10
T-16,UseCase CreateRoom + tests,usecase,3,Dev,Must,T-09
T-17,UseCase JoinRoom + tests,usecase,3,Dev,Must,T-09
T-18,UseCase CloseRoom LeaveRoom + tests,usecase,3,Dev,Must,T-09
T-19,UseCase ToggleCamera ToggleMic + tests,usecase,2,Dev,Must,T-10
T-20,UseCase NegotiateWebRTC + tests,usecase,3,Senior,Must,T-10
T-21,UseCase SendChatMessage GetChatHistory + tests,usecase,3,Dev,Should,T-09
T-22,UseCase StartStopScreenShare + tests,usecase,3,Dev,Should,T-09
T-23,HTTP Handler CreateRoom + JWT middleware,adapter,3,Dev,Must,T-16
T-24,HTTP Handler JoinRoom + GetChatHistory,adapter,3,Dev,Must,T-17
T-25,WS Handler signaling router,adapter,5,Senior,Must,T-15
T-26,HTTP Handler CloseRoom LeaveRoom,adapter,2,Dev,Should,T-18
T-27,Frontend WebRTC client create join video audio,frontend,8,Frontend,Must,T-25
T-28,Frontend chat UI,frontend,3,Frontend,Should,T-27
T-29,Frontend screen share UI,frontend,3,Frontend,Should,T-27
```
