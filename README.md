# Subscribe Tracker

Сервис для учёта подписок (MVP): регистрация/авторизация, личный кабинет, список подписок, создание/редактирование/удаление.

## План MVP
1. Фронт личного кабинета: список подписок, форма добавления/редактирования, удаление.
2. Backend авторизации: регистрация, логин, JWT.
3. Backend подписок: CRUD подписок по токену.
4. Docker Compose для локального запуска одной командой.
5. Деплой на Render: отдельные сервисы фронта и API.

## Архитектура и структура
- `frontend/` — Vite + React + TypeScript.
- `backend/` — Go API (chi + pgx), миграции в `backend/migrations/`.
- `docker-compose.yml` — локальный запуск фронта, API и Postgres.

## Схема БД (Postgres)
- `users`: id (uuid), name, email (unique), password_hash, created_at
- `subscriptions`: id (uuid), user_id (FK), service_name, bank_name, card_last4, billing_cycle (monthly/yearly), charge_date, created_at, updated_at

## API (пример)
- `POST /api/auth/register` — регистрация
- `POST /api/auth/login` — вход
- `GET /api/subscriptions` — список
- `POST /api/subscriptions` — создать
- `PUT /api/subscriptions/{id}` — обновить
- `DELETE /api/subscriptions/{id}` — удалить

## Локальный запуск (Docker Compose)
```bash
docker compose up --build
```
После старта:
- Фронт: http://localhost:5173
- API: http://localhost:8080/api

## Локальный запуск без Docker
```bash
# backend
cd backend
export DATABASE_URL=postgres://subscribe:subscribe@localhost:5432/subscribe_tracker?sslmode=disable
export JWT_SECRET=dev-secret
export CORS_ORIGINS=http://localhost:5173
go run ./cmd/server

# frontend
cd frontend
cp .env.example .env
npm install
npm run dev
```

## Деплой на Render (free)
1. Добавь репозиторий в Render.
2. Используй `render.yaml` для автоматического создания сервисов.
3. Проверь переменные окружения:
   - `DATABASE_URL`, `JWT_SECRET`, `CORS_ORIGINS` для API
   - `VITE_API_URL` для фронта

## Примечания
- Для статичного фронта используется `/api`-прокси в nginx (см. `frontend/nginx.conf`).
- При изменении API обнови типы в `frontend/src/lib/api.ts`.
