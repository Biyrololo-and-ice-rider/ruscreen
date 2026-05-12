# 🚀 Как закоммитить репозиторий из зипки

Пошаговая инструкция: от скачанного архива до первого пуша на GitHub.

---

## Шаг 1 — Распаковать и перейти в папку

```bash
unzip meet.zip -d meet
cd meet
```

---

## Шаг 2 — Инициализировать git-репозиторий

```bash
git init
git branch -M main
```

---

## Шаг 3 — Создать репозиторий на GitHub

Через веб или через GitHub CLI:

```bash
# Через gh CLI (если установлен)
gh repo create meet --private --source=. --remote=origin

# Или вручную на github.com → New repository → "meet"
# Затем добавить remote:
git remote add origin git@github.com:<your-username>/meet.git
```

---

## Шаг 4 — Первый коммит (foundation)

```bash
git add .
git commit -m "chore: initial project scaffold from design pipeline"
```

---

## Шаг 5 — Создать ветку develop и запушить обе

```bash
git push -u origin main

git checkout -b develop
git push -u origin develop
```

---

## Шаг 6 — Защитить ветки на GitHub

Settings → Branches → Add branch protection rule:

**Для `main` и `develop`:**
- ✅ Require a pull request before merging
- ✅ Require status checks to pass (CI: lint, test-unit, test-integration, build)
- ✅ Do not allow bypassing the above settings

---

## Шаг 7 — Дальнейшая работа по Git Flow

Каждая задача из [sprint backlog](docs/architecture/06-sprint-backlog.md) — отдельная ветка:

```bash
# Начало задачи T-04
git checkout develop
git pull origin develop
git checkout -b feature/us-01-create-room

# ... пишешь код ...

git add .
git commit -m "feat(domain): add Room entity with CanJoin and Close invariants"

# Когда готово — пуш и PR в develop
git push -u origin feature/us-01-create-room
# → открываешь PR на GitHub: feature/us-01-create-room → develop
```

---

## Рекомендуемый порядок первых PR (по зависимостям спринта)

| PR | Ветка | Задачи |
|---|---|---|
| PR-01 | `feature/infra-setup` | T-01: go.mod, структура, Docker |
| PR-02 | `feature/ci-lint` | T-02, T-03: golangci-lint, GitHub Actions |
| PR-03 | `feature/domain-entities` | T-04–T-08: Room, Participant, ChatMessage, тесты |
| PR-04 | `feature/ports` | T-09, T-10: интерфейсы портов |
| PR-05 | `feature/infra-token` | T-11: JWT TokenGenerator |
| PR-06 | `feature/infra-redis` | T-12–T-14: Redis репозитории + integration тесты |
| PR-07 | `feature/infra-ws-hub` | T-15: WebSocket Hub |
| PR-08 | `feature/usecase-room` | T-16–T-20: все Use Cases комнаты + WebRTC |
| PR-09 | `feature/adapters-http` | T-23–T-25: HTTP + WS handlers |
| PR-10 | `feature/frontend` | T-27: базовый WebRTC клиент |

---

## go.sum

После клонирования нужно скачать зависимости:

```bash
go mod tidy
```

Это создаст `go.sum` (он не включён в архив намеренно — генерируется автоматически).

---

## Переменные окружения

Создай `.env` (в `.gitignore`) для локальной разработки:

```env
PORT=8080
REDIS_URL=redis://localhost:6379
JWT_SECRET=local-dev-secret-change-in-prod
```

Или экспортируй напрямую:

```bash
export JWT_SECRET=local-dev-secret
make dev
```
