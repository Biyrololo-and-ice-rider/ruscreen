# 🛠 Этап 5 — DevOps Workflow

> Роль: DevOps + Senior Developer  
> Входные данные: [System Design + стек](04-system-design.md)  
> Выходные данные: Git Flow, Commits, CI, линтер → [Этап 6](06-sprint-backlog.md)

---

## Git Flow

Упрощённый Git Flow для команды 1–3 разработчика.

### Ветки

| Ветка | Назначение | Защищена | Создаётся от |
|---|---|---|---|
| `main` | Продакшн, всегда стабильна | ✅ | — |
| `develop` | Интеграционная, база для фич | ✅ | — |
| `feature/<us-id>-<slug>` | Новая фича | ❌ | develop |
| `fix/<slug>` | Баг-фикс | ❌ | develop / main |
| `release/<version>` | Подготовка релиза | ❌ | develop |
| `hotfix/<slug>` | Срочный прод-фикс | ❌ | main |

### Правила merge

```
feature/* → develop        squash merge (через PR)
fix/*     → develop        squash merge (через PR)
release/* → main + develop merge commit
hotfix/*  → main + develop merge commit
```

### Примеры имён веток

```
feature/us-01-create-room
feature/us-05-webrtc-negotiation
feature/us-08-chat
fix/ws-disconnect-during-sdp
hotfix/room-full-boundary
```

---

## Conventional Commits

### Формат

```
<type>(<scope>): <краткое описание>

[опциональное тело]

[BREAKING CHANGE: ... / Closes #issue]
```

### Типы

| Тип | Когда |
|---|---|
| `feat` | Новая функциональность |
| `fix` | Исправление бага |
| `refactor` | Рефакторинг без изменения поведения |
| `test` | Добавление/изменение тестов |
| `docs` | Только документация |
| `chore` | Зависимости, конфиги |
| `ci` | GitHub Actions |
| `build` | Dockerfile, Makefile |
| `perf` | Оптимизация производительности |

### Scopes проекта

```
domain, usecase, adapter, infra, ws, redis, token, http, room, chat, signaling
```

### Примеры коммитов

```
feat(usecase): implement JoinRoom with room capacity check
feat(ws): add WebSocket hub with goroutine fan-out
feat(infra): implement Redis RoomRepository with TTL
fix(infra): handle Redis key expiry in RoomRepository.FindByID
refactor(domain): extract RoomStatus as typed iota const
test(usecase): add mock-based unit tests for CreateRoom
chore(deps): upgrade go-redis to v9.5.1
ci: add integration test step with Redis service container
build: add multi-stage Dockerfile for lean production image
docs(readme): add API endpoint table and quick start guide
```

---

## .golangci.yml

```yaml
run:
  timeout: 5m
  go: "1.22"

linters:
  enable:
    - errcheck       # проверка необработанных ошибок
    - gosimple       # упрощение кода
    - govet          # анализ потенциальных ошибок
    - ineffassign    # неиспользуемые присвоения
    - staticcheck    # статический анализ
    - unused         # неиспользуемый код
    - goimports      # форматирование импортов
    - misspell       # опечатки в комментариях
    - revive         # расширенный линтер (замена golint)
    - wrapcheck      # проверка оборачивания ошибок
    - depguard       # контроль импортов по слоям

linters-settings:
  depguard:
    rules:
      domain-isolation:
        files: ["**/domain/**"]
        deny:
          - pkg: "github.com/yourorg/meet/internal/infra"
            desc: "domain не должен импортировать infra"
          - pkg: "github.com/yourorg/meet/internal/adapter"
            desc: "domain не должен импортировать adapter"
      usecase-isolation:
        files: ["**/usecase/**"]
        deny:
          - pkg: "github.com/yourorg/meet/internal/infra"
            desc: "usecase использует порты, а не infra напрямую"
```

---

## GitHub Actions CI

```yaml
# .github/workflows/ci.yml
name: CI

on:
  pull_request:
    branches: [main, develop]
  push:
    branches: [develop]

jobs:
  lint:
    name: Lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"
          cache: true
      - uses: golangci/golangci-lint-action@v4
        with:
          version: latest

  test-unit:
    name: Unit Tests
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"
          cache: true
      - run: make test-unit

  test-integration:
    name: Integration Tests
    runs-on: ubuntu-latest
    services:
      redis:
        image: redis:7-alpine
        ports: ["6379:6379"]
        options: >-
          --health-cmd "redis-cli ping"
          --health-interval 5s
          --health-timeout 3s
          --health-retries 5
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"
          cache: true
      - run: make test-integration
        env:
          REDIS_URL: redis://localhost:6379

  build:
    name: Build
    runs-on: ubuntu-latest
    needs: [lint, test-unit, test-integration]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"
          cache: true
      - run: make build
      - run: docker build -t meet:ci .
```

---

## Makefile

```makefile
.PHONY: build run dev test test-unit test-integration lint docker-up docker-down generate

build:
	go build -o bin/meet ./cmd/server

run: build
	./bin/meet

dev:
	docker-compose up -d && go run ./cmd/server

test-unit:
	go test ./internal/domain/... ./internal/usecase/... -v -race -count=1

test-integration:
	go test ./internal/infra/... ./internal/adapter/... \
		-v -race -count=1 -tags=integration

test: test-unit test-integration

lint:
	golangci-lint run ./...

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

generate:
	go generate ./...
```

---

## PR Checklist

```markdown
## Чеклист перед merge

### Код
- [ ] `make build` — компилируется без ошибок
- [ ] `make lint` — golangci-lint чистый
- [ ] `make test` — все тесты зелёные
- [ ] Покрытие нового кода ≥ 80% (domain + usecase)

### Архитектура
- [ ] domain/ не импортирует infra/ или adapter/
- [ ] usecase/ не обращается к infra/ напрямую (только через порты)
- [ ] adapter/ использует только DTO, не domain.Entity в HTTP-ответах

### Качество
- [ ] Все публичные типы и методы задокументированы
- [ ] Ошибки обёрнуты через fmt.Errorf("context: %w", err)
- [ ] Нет прямого panic в бизнес-логике

### PR
- [ ] Название в формате Conventional Commits
- [ ] Описание содержит ссылку на задачу
- [ ] Минимум 1 аппрув перед merge
```

---

## Definition of Done

```
✅ Код компилируется, go vet чистый
✅ Unit-тесты написаны, проходят (покрытие ≥ 80% для domain + usecase)
✅ golangci-lint без ошибок
✅ PR прошёл CI (lint + test-unit + test-integration + build)
✅ Минимум 1 аппрув
✅ Смержено в develop через squash merge
✅ Задача закрыта в трекере
```

---

→ Следующий этап: [06-sprint-backlog.md](06-sprint-backlog.md)
