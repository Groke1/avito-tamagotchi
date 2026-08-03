# 🐥 Тамагочи Авито — Клиентское веб-приложение

Игровая система лояльности и мотивации для пользователей Авито. Веб-клиент разработан в рамках хакатона на современном технологическом стеке с использованием архитектуры **Feature-Sliced Design (FSD)**.

---

## 🛠️ Технологический стек

- **Ядро:** React 19, TypeScript
- **Сборка:** Vite 8 (нативная поддержка TS-путей)
- **Архитектура:** Feature-Sliced Design (FSD)
- **Стилизация:** Tailwind CSS v4 (`@tailwindcss/vite`), shadcn/ui (Slate / New York)
- **Стейт-менеджмент & API:** Redux Toolkit 2, RTK Query (`@reduxjs/toolkit/query/react`)
- **Маршрутизация:** React Router DOM v7
- **Иконки & Утилиты:** `lucide-react`, `clsx`, `tailwind-merge`, `class-variance-authority`
- **Качество кода:** ESLint, Prettier

---

## 📐 Архитектура проекта (Feature-Sliced Design)

Исходный код расположен в `src/` и структурирован строго по слоям FSD:

```
src/
├── app/            # Инициализация приложения (Providers, Store, Router, Глобальные стили)
├── pages/          # Компоненты страниц (HomePage, PetPage, LeaderboardPage)
├── widgets/        # Крупные блоки интерфейса (Header, PetCardWidget, DailySummary)
├── features/       # Действия пользователя (feed-pet, complete-task, claim-reward)
├── entities/       # Бизнес-сущности (pet, user, task, reward, leaderboard)
└── shared/         # Переиспользуемый код без бизнес-логики
    ├── api/        # Базовый RTK Query slice
    ├── config/     # Константы и конфигурация
    ├── lib/        # Утилиты (cn helper, formatters)
    ├── types/      # Глобальные TS типы
    └── ui/         # Атомарные UI-компоненты (shadcn/ui: Button, Card, Dialog и т.д.)
```

---

## 🚀 Команды для локальной разработки

Установка зависимостей:

```bash
npm install
```

Запуск сервера разработки:

```bash
npm run dev
```

Проверка типов и сборка production-бандла:

```bash
npm run build
```

Проверка кода линтером:

```bash
npm run lint
```

Автоматическое форматирование кода и сортировка импортов:

```bash
npm run format
```

---

## 📚 Документация и Контракты

- Контракт API: [`docs/api-spec.yaml`](../docs/api-spec.yaml)
- Регламент Git Flow: [`docs/GIT_FLOW.md`](../docs/GIT_FLOW.md)
- Инструкции для AI-агентов: [`AGENTS.md`](./AGENTS.md)
