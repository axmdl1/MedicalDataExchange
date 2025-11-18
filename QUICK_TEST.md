# ⚡ Быстрое тестирование системы

## 🚀 Запуск (5 минут)

### 1. Запустите Docker

```bash
cd /Users/damir/clinic-data/MedicalDataExchange

# Запустить все сервисы
docker-compose up -d

# Подождать 10 секунд
sleep 10

# Проверить статус
docker-compose ps
```

**Должны быть запущены:**
- ✅ mde-postgres
- ✅ mde-core
- ✅ mde-dt
- ✅ mde-user
- ✅ mde-frontend

---

### 2. Откройте фронтенд

```
Браузер → http://localhost:8080
```

---

## 🧪 Быстрый тест (2 минуты)

### Вариант A: Через фронтенд

**Шаг 1:** Откройте новую страницу
```
http://localhost:8080/patient-access.html
```

**Шаг 2:** Вы увидите 4 вкладки:
- 📨 Создать запрос
- 📋 Мои запросы
- ⏳ Ожидают одобрения
- 👁️ Просмотр данных

**Готово!** Страница работает.

---

### Вариант B: Через API (curl)

```bash
# 1. Создать запрос на доступ
curl -X POST http://localhost:8099/patient-access/request \
  -H "Content-Type: application/json" \
  -d '{
    "patient_id": 1,
    "clinic_id": 1,
    "medical_data_id": 1
  }'

# Ожидаемый ответ:
# {
#   "request_id": 1,
#   "status": "pending",
#   "message": "Запрос на доступ создан..."
# }

# 2. Одобрить запрос
curl -X POST http://localhost:8099/patient-access/approve/1 \
  -H "Content-Type: application/json" \
  -d '{
    "approver_id": 2
  }'

# Ожидаемый ответ:
# {
#   "request_id": 1,
#   "access_token": "mock-token-abc123xyz",
#   "expires_at": "2024-01-01T12:15:00Z",
#   "status": "approved"
# }

# 3. Получить данные по токену
curl "http://localhost:8099/patient-access/data?token=mock-token-abc123xyz"

# Ожидаемый ответ:
# {
#   "diagnosis": "ОРВИ",
#   "complaint": "Температура...",
#   ...
# }
```

---

## 📊 Проверка в БД (1 минута)

```bash
# Подключитесь к PostgreSQL
docker exec -it mde-postgres psql -U postgres -d medical

# Проверьте таблицы
\dt patient_access_requests
\dt temporary_patient_data

# Посмотрите данные
SELECT * FROM patient_access_requests LIMIT 5;
SELECT * FROM temporary_patient_data LIMIT 5;

# Выход
\q
```

---

## 📋 Полный сценарий тестирования

Смотрите подробный гайд:
```
TESTING_SCENARIOS.md
```

---

## 🎯 Ключевые URL

| Что | URL |
|-----|-----|
| Главная | http://localhost:8080 |
| **Временный доступ** | **http://localhost:8080/patient-access.html** |
| API | http://localhost:8099 |
| PostgreSQL | localhost:5442 |

---

## 🔑 Тестовые аккаунты (если есть seed данные)

| Роль | Email | Password |
|------|-------|----------|
| Пациент | patient@example.com | password123 |
| Врач | doctor@clinic.com | password123 |
| Админ | admin@example.com | password123 |

---

## ✅ Что проверить

1. **Фронтенд:**
   - [ ] Открывается http://localhost:8080/patient-access.html
   - [ ] 4 вкладки отображаются
   - [ ] Можно переключаться между табами

2. **API:**
   - [ ] POST /patient-access/request работает
   - [ ] POST /patient-access/approve/{id} работает
   - [ ] GET /patient-access/data?token=xxx работает

3. **База данных:**
   - [ ] Таблица patient_access_requests существует
   - [ ] Таблица temporary_patient_data существует

---

## 🐛 Если что-то не работает

### Проблема: "Cannot connect to Docker daemon"
```bash
# Запустите Docker Desktop
open -a Docker
```

### Проблема: "Connection refused на порту 8099"
```bash
# Проверьте логи
docker logs mde-core

# Перезапустите
docker-compose restart core-service
```

### Проблема: Страница 404
```bash
# Убедитесь что файл существует
ls frontend/patient-access.html

# Перезапустите frontend
docker-compose restart frontend
```

### Проблема: Пустые списки на фронте
```bash
# Проверьте консоль браузера (F12)
# Посмотрите вкладку Network на ошибки API
```

---

## 📝 Краткое резюме

**Что было добавлено:**
1. ✅ REST API endpoints (6 штук)
2. ✅ Фронтенд страница с 4 вкладками
3. ✅ API клиент (api.js обновлен)
4. ✅ Новые таблицы в БД
5. ✅ Cleanup scheduler (фоновая задача)

**Что работает сейчас:**
- ✅ Создание запросов на доступ
- ✅ Просмотр списка запросов
- ✅ Одобрение запросов клиникой
- ✅ Просмотр данных по токену
- ✅ Обратный отсчет 15 минут
- ✅ Mock данные возвращаются

**Что нужно для полной работы:**
- ⏳ Hyperledger Fabric (опционально)
- ⏳ Интеграция handlers с data-transfer-service
- ⏳ Реальное шифрование данных

---

## 🎉 Готово!

Система готова к демонстрации!

Можете показать:
- Создание запроса пациентом
- Одобрение клиникой
- Временный доступ (15 минут)
- Обратный отсчет
- Автоматическое удаление

**Хотите протестировать?**

Следуйте инструкциям в `TESTING_SCENARIOS.md`
