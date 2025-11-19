# 🏥 Medical Data Exchange - Система управления медицинскими данными

## 🎯 Что это?

Полноценная система управления медицинскими данными с **реальной** функциональностью:

- ✅ **Временный доступ пациента** - данные доступны только 10-15 минут
- ✅ **Данные хранятся в клинике** - пациент не может хранить данные у себя
- ✅ **Автоматическое удаление** - через 15 минут данные удаляются автоматически
- ✅ **Блокчейн аудит** - Hyperledger Fabric для прозрачности
- ✅ **Реальные медицинские данные** - не моки, а полноценные клинические записи

## 🚀 Быстрый старт

```bash
# 1. Запуск
docker-compose up -d

# 2. Тест
./test-flow.sh

# Готово! Система работает на http://localhost:8080
```

## 📋 Как это работает?

### 1️⃣ Пациент создает запрос
```
Пациент: "Хочу посмотреть свои анализы"
  ↓
POST /api/patient-access/request
  ↓
Запрос сохранен (status: pending)
```

### 2️⃣ Клиника одобряет
```
Клиника: "Одобряю доступ"
  ↓
POST /api/patient-access/approve/1
  ↓
ЧТО ПРОИСХОДИТ:
  1. Данные извлекаются из БД клиники
  2. Шифруются AES-256-GCM
  3. Создается временная копия
  4. Устанавливается таймер на 15 минут
  5. Генерируется уникальный токен
  ↓
Пациент получает токен доступа
```

### 3️⃣ Пациент получает данные
```
GET /api/patient-access/data?token=abc123...
  ↓
Система проверяет:
  ✓ Токен валидный?
  ✓ Срок не истек? (< 15 мин)
  ✓ Не отозван?
  ↓
Расшифровывает и отдает данные
```

### 4️⃣ Автоматическое удаление
```
Через 15 минут:
  ↓
Scheduler (каждые 5 минут):
  ↓
DELETE FROM temporary_patient_data
WHERE expires_at < NOW()
  ↓
Временные данные УДАЛЕНЫ
Оригинал остался в клинике ✓
```

## 🏗️ Архитектура

```
┌────────────────────────────────┐
│    КЛИНИКА (База данных)       │
│                                 │
│  medical_data:                  │
│  • Все данные пациента          │
│  • НИКОГДА не удаляются         │
│  • 50+ полей                    │
│                                 │
│  temporary_patient_data:        │
│  • Зашифрованная копия          │
│  • Срок: 15 минут               │
│  • АВТОМАТИЧЕСКИ удаляется      │
└────────────────────────────────┘
        ▲              │
        │              │
    Запрос        Одобрение
        │              │
        │              ▼
┌───────────┐    ┌──────────┐
│  Пациент  │←───│  Клиника │
│ (ничего   │    │(постоянное│
│  не       │    │ хранение) │
│  хранит)  │    └──────────┘
└───────────┘
     │
     │ Через 15 мин
     ▼
   ❌ EXPIRED
```

## 💾 Примеры данных

### Реальные медицинские записи:

**ОРВИ (Patient 1):**
```json
{
  "visit_date": "2025-11-16",
  "department": "General Medicine",
  "chief_complaint": "High fever and severe headache",
  "temperature": 38.7,
  "blood_pressure": "125/78",
  "heart_rate": 92,
  "diagnosis": "Acute upper respiratory tract infection",
  "diagnosis_code": "J06.9",
  "prescribed_medications": [
    {"name": "Paracetamol", "dosage": "500mg", "frequency": "Every 6 hours"}
  ],
  "lab_results": {
    "wbc": 8.5,
    "hemoglobin": 13.2,
    "platelets": 245
  },
  "doctor_notes": "Patient presents with classic viral URTI symptoms..."
}
```

## 🔐 Безопасность

### Шифрование
- **AES-256-GCM** - военный стандарт шифрования
- **SHA-256** - проверка целостности данных
- **Уникальные токены** - невозможно угадать

### Блокчейн
- **Hyperledger Fabric** - immutable audit trail
- Все действия записываются
- Невозможно изменить историю
- Прозрачность для аудита

## 📊 Что включено?

| Компонент | Описание | Статус |
|-----------|----------|--------|
| **Медицинские данные** | 50+ полей: диагнозы, анализы, назначения | ✅ |
| **Система запросов** | Patient → Clinic → Approve/Reject | ✅ |
| **Временный доступ** | Ровно 15 минут, потом удаление | ✅ |
| **Шифрование** | AES-256-GCM для временных данных | ✅ |
| **Автоочистка** | Scheduler каждые 5 минут | ✅ |
| **Блокчейн** | Hyperledger Fabric audit trail | ✅ |
| **API** | RESTful endpoints | ✅ |
| **Документация** | Полная на русском и английском | ✅ |

## 🧪 Тестирование

### Автоматический тест
```bash
./test-flow.sh
```

### Ручное тестирование

**1. Создать запрос:**
```bash
curl -X POST http://localhost:8099/api/patient-access/request \
  -H "Content-Type: application/json" \
  -d '{"patient_id": 1, "clinic_id": 1, "medical_data_id": 1}'
```

**2. Одобрить:**
```bash
curl -X POST http://localhost:8099/api/patient-access/approve/1 \
  -H "Content-Type: application/json" \
  -d '{"approver_id": 10}'
```

**3. Получить данные:**
```bash
curl "http://localhost:8099/api/patient-access/data?token=ТОКЕН_ИЗ_ШАГА_2"
```

**4. Подождать 15 минут и повторить шаг 3:**
```
Ожидается: "error": "access token has expired"
```

## 📚 Документация

- **START_HERE.md** - Начните отсюда
- **IMPLEMENTATION_GUIDE.md** - Полное техническое руководство
- **FINAL_SUMMARY.md** - Итоговая сводка
- **WHATS_NEW.md** - Что нового

## 🎯 Ключевые особенности

### ✅ Не моки - реальная система
- Настоящее шифрование (AES-256-GCM)
- Реальный scheduler (каждые 5 минут)
- Настоящая БД с индексами
- Полноценный блокчейн (Hyperledger Fabric)

### ✅ Данные в клинике
- `medical_data` - постоянное хранение
- `temporary_patient_data` - временная копия
- Оригинал НИКОГДА не удаляется
- Пациент не может хранить у себя

### ✅ Автоматическое удаление
- Scheduler запускается каждые 5 минут
- Проверяет `expires_at < NOW()`
- Удаляет истекшие данные
- Обновляет статусы

### ✅ Production-ready
- Микросервисная архитектура
- gRPC между сервисами
- Docker-compose для запуска
- Graceful shutdown
- Логирование и мониторинг

## 🛠️ Технологии

- **Backend:** Go (Golang)
- **Database:** PostgreSQL
- **Blockchain:** Hyperledger Fabric
- **Encryption:** AES-256-GCM
- **Hash:** SHA-256
- **API:** REST + gRPC
- **Containerization:** Docker

## 📞 Поддержка

Если что-то не работает:

1. **Проверьте логи:**
   ```bash
   docker logs mde-dt -f
   docker logs mde-postgres
   ```

2. **Проверьте БД:**
   ```bash
   docker exec -it mde-postgres psql -U postgres -d medical
   ```

3. **Читайте документацию:**
   - `START_HERE.md` для quick start
   - `IMPLEMENTATION_GUIDE.md` для troubleshooting

## 🎉 Готово к использованию!

Система **полностью рабочая** и готова к запуску.

**Дальше:**
1. `docker-compose up -d` - запуск
2. `./test-flow.sh` - тест
3. Наслаждайтесь! 🚀

---

**Made with ❤️ for real medical data management**
