# Microservice Sneaker Store

Репозиторий: https://github.com/Abegerov4/microservice-sneaker-store

## Структура проекта

| Сервис | Папка | Ответственный |
|--------|-------|---------------|
| User Service | `user-service/` | — |
| Order Service | `order-service/` | — |
| Product Service | `product-service/` | — |

---

## Инструкция для каждого участника

### Шаг 1 — Установить Git

Скачать: https://git-scm.com/downloads  
Проверить установку:
```bash
git --version
```

Настроить своё имя и email (один раз):
```bash
git config --global user.name "Твоё Имя"
git config --global user.email "твой@email.com"
```

---

### Шаг 2 — Принять приглашение

1. Зайди на почту, которую указал при регистрации на GitHub
2. Найди письмо от GitHub с темой "You've been invited to..."
3. Нажми **Accept invitation**

---

### Шаг 3 — Клонировать репозиторий

```bash
git clone https://github.com/Abegerov4/microservice-sneaker-store.git
cd microservice-sneaker-store
```

---

### Шаг 4 — Создать свою ветку

> Каждый работает в своей ветке. В `main` напрямую не коммитим.

```bash
# Для user-service:
git checkout -b feature/user-service

# Для order-service:
git checkout -b feature/order-service

# Для product-service:
git checkout -b feature/product-service
```

---

### Шаг 5 — Создать папку своего сервиса

```bash
# Например для user-service:
mkdir user-service
cd user-service
# ... создавай свои файлы ...
```

---

### Шаг 6 — Ежедневная работа (коммиты)

```bash
# 1. Посмотреть что изменилось
git status

# 2. Добавить файлы в коммит
git add .
# или конкретный файл:
git add user-service/index.js

# 3. Создать коммит с описанием
git commit -m "Add user registration endpoint"

# 4. Отправить на GitHub
git push origin feature/user-service
```

**Правила для сообщений коммита:**
- `Add ...` — добавил новую функцию
- `Fix ...` — исправил баг
- `Update ...` — обновил что-то существующее
- `Remove ...` — удалил что-то

---

### Шаг 7 — Получить последние изменения от команды

Делай это каждый раз перед началом работы:

```bash
# Переключиться на main
git checkout main

# Скачать последние изменения
git pull origin main

# Вернуться на свою ветку
git checkout feature/user-service   # (или свою ветку)

# Влить изменения из main в свою ветку
git merge main
```

---

### Шаг 8 — Создать Pull Request (когда часть готова)

1. Зайди на https://github.com/Abegerov4/microservice-sneaker-store
2. GitHub покажет жёлтую кнопку **"Compare & pull request"** — нажми её
3. Напиши что сделал в описании
4. Нажми **"Create pull request"**
5. Напиши остальным участникам — пусть проверят и approve-ят

---

## Частые ошибки

**Ошибка: `rejected - non-fast-forward`**
```bash
git pull origin feature/user-service --rebase
git push origin feature/user-service
```

**Ошибка: `Please tell me who you are`**
```bash
git config --global user.name "Твоё Имя"
git config --global user.email "твой@email.com"
```

**Случайно закоммитил в main:**
```bash
git checkout -b feature/user-service   # создай ветку из текущего состояния
git checkout main
git reset --hard origin/main           # сбросить main обратно
```

---

## Структура веток

```
main
├── feature/user-service
├── feature/order-service
└── feature/product-service
```

---

## Контакты команды

| Участник | Сервис | GitHub |
|----------|--------|--------|
| Aldiyar | владелец | @Abegerov4 |
| — | user-service | — |
| — | order-service | — |
| — | product-service | — |
