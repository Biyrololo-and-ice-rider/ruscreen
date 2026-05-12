# Contributing to Meet

Привет! Это гайд для всех, кто работает над проектом.  
Здесь — всё, что нужно знать, чтобы писать код единообразно и без конфликтов.

---

## Быстрый старт

```bash
git clone <repo-url>
cd meet
make docker-up    # поднять Redis
make run          # запустить сервер на :8080
make test         # убедиться, что всё работает
make lint         # проверить линтер
```

---

## Git Flow

Используем упрощённый Git Flow для небольшой команды.

### Ветки

| Ветка | Назначение | Защищена |
|---|---|---|
| `main` | Продакшн, всегда стабильна | ✅ |
| `develop` | Интеграционная, база для всех фич | ✅ |
| `feature/<us-id>-<slug>` | Новая фича | ❌ |
| `fix/<slug>` | Баг-фикс | ❌ |
| `release/<version>` | Подготовка релиза | ❌ |
| `hotfix/<slug>` | Срочный прод-фикс | ❌ |

**Примеры имён веток:**
```
feature/us-01-create-room
feature/us-05-webrtc-negotiation
fix/ws-disconnect-sdp
hotfix/room-full-check
```

### Правила merge

```
feature/* → develop        squash merge (через PR)
fix/*     → develop        squash merge (через PR)
release/* → main + develop merge commit
hotfix/*  → main + develop merge commit
```

> Прямые коммиты в `main` и `develop` — **запрещены**. Только через PR.

---

## Conventional Commits

Формат:
```
<type>(<scope>): <описание>

[опциональное тело]

[опциональный footer: BREAKING CHANGE / Closes #issue]
```

### Типы

| Тип | Когда |
|---|---|
| `feat` | Новая функциональность |
| `fix` | Исправление бага |
| `refactor` | Рефакторинг без изменения поведения |
| `test` | Тесты |
| `docs` | Только документация |
| `chore` | Зависимости, конфиги |
| `ci` | GitHub Actions |
| `build` | Dockerfile, Makefile |
| `perf` | Оптимизация |

### Scopes проекта

```
domain, usecase, adapter, infra, ws, redis, token, http, room, chat, signaling
```

### Примеры

```
feat(usecase): implement JoinRoom interactor with capacity check
fix(infra): handle Redis key expiry in RoomRepository.FindByID
feat(ws): add WebSocket hub with broadcast fan-out
refactor(domain): extract RoomStatus as typed iota
test(usecase): add mock-based unit tests for CreateRoom
chore(ci): add integration test step with Redis service
docs(readme): add API endpoint table
```

---

## Правила именования

- **Файлы:** `snake_case.go` — один файл = один Use Case или одна Entity
- **Пакеты:** строчные, без подчёркиваний (`domain`, `usecase`, `redis`)
- **Интерфейсы:** `PascalCase` с суффиксом по типу (`RoomRepository`, `SignalingGateway`)
- **Ошибки:** `var ErrXxx = errors.New(...)` в `domain/errors.go`
- **Use Cases:** файл = `create_room.go`, структура = `CreateRoomUseCase`

---

## Архитектурные правила

```
domain/ → только stdlib
usecase/ → только domain/ + usecase/port/ + stdlib
adapter/ → usecase/ + domain/ (только Request/Response Models)
infra/ → всё что угодно (Redis, JWT, WS)
```

**Нарушения, которых быть не должно:**
- `domain/` импортирует `infra/` или `adapter/` → **запрещено**
- `usecase/` вызывает Redis напрямую → **запрещено**, только через порт
- `adapter/` возвращает domain.Entity в HTTP ответе → **запрещено**, только DTO

Эти правила проверяются линтером (`depguard` в `.golangci.yml`).

---

## PR Checklist

См. шаблон: [.github/pull_request_template.md](.github/pull_request_template.md)

Коротко: перед открытием PR убедись:
```
make build   # компилируется
make test    # тесты зелёные
make lint    # линтер доволен
```

---

## Definition of Done

Задача считается выполненной когда:

- ✅ Код компилируется без ошибок
- ✅ Unit-тесты написаны, проходят (покрытие нового кода ≥ 80%)
- ✅ `golangci-lint` без ошибок
- ✅ PR прошёл CI (lint + test + build)
- ✅ Минимум 1 код-ревью аппрув
- ✅ Смержено в `develop`
- ✅ Задача закрыта в трекере
