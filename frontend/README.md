# Тамагочи Авито — Клиентское веб-приложение (Frontend)

Игровой веб-клиент системы лояльности и мотивации для пользователей Авито. Приложение разработано в рамках хакатона на современном технологическом стеке с использованием архитектуры **Feature-Sliced Design (FSD)**.

---

## 🛠️ Технологический стек Frontend

- **Ядро:** React 19, TypeScript
- **Сборка:** Vite 8 (нативная поддержка TS-алиасов `@/...`)
- **Архитектура:** Feature-Sliced Design (FSD)
- **Стилизация:** Tailwind CSS v4 (`@tailwindcss/vite`), CSS Variable Design Tokens
- **Стейт-менеджмент & API:** Redux Toolkit 2, RTK Query (`@reduxjs/toolkit/query/react`)
- **Маршрутизация:** React Router DOM v7 (с авторизационными гардами)
- **Тестирование:** Vitest, React Testing Library, `@testing-library/jest-dom`, JSDOM
- **Иконки & Утилиты:** `lucide-react`, `clsx`, `tailwind-merge`, `class-variance-authority`
- **Качество кода:** ESLint, Prettier

---

## 📐 Архитектура проекта (Feature-Sliced Design)

Исходный код расположен в `src/` и строго разделен по слоям архитектуры FSD:

```
src/
├── app/            # Инициализация (Providers, Redux Store, Router)
├── pages/          # Страницы приложения (HomePage, DailyReportPage, TasksPage, RewardsPage, LeaderboardPage, AuthPages)
├── widgets/        # Крупные блоки интерфейса (Header, PetDashboard, DailyReport, Leaderboard, TasksList, RewardsList)
├── features/       # Интерактивные сценарии (pet-actions, auth)
├── entities/       # Бизнес-сущности и модель данных (pet, user, task, reward, leaderboard)
└── shared/         # Независимый переиспользуемый код
    ├── api/        # Базовый RTK Query slice и WebSocket клиенты
    ├── config/     # Маршруты, темы и константы
    ├── lib/        # Утилиты, форматирование дат, валидаторы, guards
    └── ui/         # Атомарные UI-компоненты (Button, Input, Card, Skeleton, Dialog, Tabs)
```

## ✨ Особенности реализации Frontend

1. **Адаптивный Mobile-First дизайн:**
   - Полная отзывчивость всех страниц и виджетов (мобильные устройства от 320px, планшеты, десктопы).

2. **WebSocket & Реальное время:**
   - Поддержка синхронизации состояния питомца (опыт, уровень, монеты, стрик) между вкладками и клиентами в реальном времени.

3. **Авторизация и защищенные маршруты:**
   - Кастомные React Router гарды (`AuthGuard`, `GuestGuard`) для предотвращения несанкционированного доступа.
   - лоадер (`rootLoader`) для проверки авторизации пользователя

4. **Интерактивные механики и анимации:**
   - Визуальные эффекты при кормлении и поглаживании питомца с анимациями, прогресс-барами.

5. **Надежное автоматическое тестирование (Vitest):**
   - Набор из тестов, покрывающих Redux-слайсы, бизнес-хуки (`useHandleFeed`, `useHandleStroke`, `useDailyReportCards`), схемы валидации Zod/Form и UI-компоненты.

---

## 🚀 Команды для локальной разработки и тестирования

### 1. Установка зависимостей:

```bash
cd frontend
npm install
```

### 2. Запуск сервера разработки:

```bash
npm run dev
```

### 3. Запуск модульных и компонентных тестов (Vitest):

```bash
npm run test:run
```

Запуск тестов в режиме наблюдения (watch mode):

```bash
npm run test
```

### 4. Проверка типов и сборка Production-бандла:

```bash
npm run build
```

### 5. Проверка линтером и форматирование:

```bash
npm run lint
npm run format
```
