# Практическое занятие №13: RabbitMQ – отправка и получение сообщений

## Описание

Реализована асинхронная коммуникация между сервисами через **RabbitMQ**.  
При создании задачи сервис `tasks` публикует событие `task.created` в очередь `task_events`.  
Отдельный worker‑процесс читает сообщения из очереди, логирует их и подтверждает обработку (ack).

**Producer** – сервис `tasks` (HTTP на порту 8082)  
**Consumer** – worker (читает очередь `task_events`)

---

## Технологии

- **Go** 1.21+
- **RabbitMQ** (брокер сообщений)
- **Docker Compose** (запуск RabbitMQ)
- **amqp091-go** – клиент для RabbitMQ

---

## Структура проекта

![Структура проекта](screen/structure.png)


---

## Запуск

### 1. Запустить RabbitMQ

```bash
cd deploy/rabbit
docker compose up -d
```
![Структура проекта](screen/guest_guest.png)
![Структура проекта](screen/page.png)

### 2. Запустить worker (consumer)
![Структура проекта](screen/worker_started.png)

### 3. Запустить сервис tasks (producer)
![Структура проекта](screen/tasks_started.png)

### Проверка работы
![Структура проекта](screen/201_created.png)

### Логи сервиса tasks
![Структура проекта](screen/tasks_logs.png)

### Логи worker
![Структура проекта](screen/worker_logs.png)

### Проверка очереди в Management UI
![Структура проекта](screen/page.png)

