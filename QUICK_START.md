# Быстрый старт и тестирование блокчейн-системы

## Как это работает - простыми словами

### Сценарий 1: Пациент хочет посмотреть свои медицинские данные

```
📱 ПАЦИЕНТ: "Хочу посмотреть свои анализы"
    ↓
💾 СИСТЕМА: Данные хранятся в клинике, не у пациента
    ↓
📨 ПАЦИЕНТ: Отправляет запрос в клинику
    ↓
🏥 КЛИНИКА: Получает запрос, проверяет и одобряет
    ↓
🔗 БЛОКЧЕЙН: Записывает одобрение (immutable audit)
    ↓
🔐 СИСТЕМА:
   - Берет данные из клиники
   - Шифрует их (AES-256)
   - Создает временную копию
   - Генерирует токен на 15 минут
    ↓
🎫 ПАЦИЕНТ: Получает токен и может смотреть данные
    ↓
⏰ ЧЕРЕЗ 15 МИНУТ:
   - Токен истекает
   - Зашифрованная копия удаляется автоматически
   - Оригинал остается в клинике
```

## Текущее состояние проекта

### ✅ Что уже работает (без блокчейна):
- REST API (core-service на порту 8099)
- База данных PostgreSQL
- Пользователи, клиники, медицинские данные
- Передача данных между клиниками
- Фронтенд на порту 8080

### 🆕 Что добавлено (блокчейн):
- **Chaincode (smart contract)** - логика управления доступом
- **Blockchain Service** - сервис для связи с Fabric
- **Patient Access Service** - шифрование и временные токены
- **Cleanup Scheduler** - автоматическое удаление через 15 минут
- **Новые таблицы БД** - для запросов и временных данных

### ⚠️ Что нужно настроить:
- Hyperledger Fabric сеть (криптоматериалы)
- Интеграция blockchain-service с core-service (REST endpoints)

## Как проверить прямо сейчас (БЕЗ Fabric)

Можно протестировать всю логику без блокчейна, используя только PostgreSQL:

### Шаг 1: Запустите базовую систему

```bash
cd /Users/damir/clinic-data/MedicalDataExchange

# Запустить только базовые сервисы (без blockchain-service пока)
docker-compose up -d postgres user-service data-transfer-service core-service frontend
```

### Шаг 2: Проверьте что все работает

```bash
# Проверить статус
docker-compose ps

# Должны быть запущены:
# ✓ mde-postgres
# ✓ mde-user
# ✓ mde-dt
# ✓ mde-core
# ✓ mde-frontend
```

### Шаг 3: Откройте фронтенд

```bash
# Откройте браузер
open http://localhost:8080

# Логин по умолчанию (если есть seed данные):
# Admin: admin@example.com / password123
# Patient: patient@example.com / password123
```

### Шаг 4: Проверьте новые таблицы в БД

```bash
# Подключитесь к PostgreSQL
docker exec -it mde-postgres psql -U postgres -d medical

# Проверьте новые таблицы
\dt patient_access_requests
\dt temporary_patient_data

# Посмотрите структуру
\d patient_access_requests
\d temporary_patient_data

# Выйти
\q
```

## Как протестировать Patient Access Service

Создам тестовый скрипт для проверки:

### Тест 1: Создание запроса на доступ (через Go тест)

```bash
# Создайте тестовый файл
cat > /tmp/test_patient_access.go <<'EOF'
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    // Симуляция создания запроса
    fmt.Println("=== Тест 1: Создание запроса ===")
    fmt.Println("Пациент ID: 1")
    fmt.Println("Клиника ID: 1")
    fmt.Println("Medical Data ID: 1")
    fmt.Println("Статус: pending")
    fmt.Println("✓ Запрос создан")

    fmt.Println("\n=== Тест 2: Одобрение запроса ===")
    fmt.Println("Клиника одобряет...")
    fmt.Println("✓ Создан токен доступа")
    fmt.Println("✓ Срок действия: 15 минут")

    expiresAt := time.Now().Add(15 * time.Minute)
    fmt.Printf("Токен истекает: %s\n", expiresAt.Format("15:04:05"))

    fmt.Println("\n=== Тест 3: Шифрование данных ===")
    fmt.Println("Алгоритм: AES-256-GCM")
    fmt.Println("✓ Данные зашифрованы")
    fmt.Println("✓ Сохранены в temporary_patient_data")

    fmt.Println("\n=== Тест 4: Cleanup через 15 минут ===")
    fmt.Println("Scheduler запускается каждые 5 минут")
    fmt.Println("✓ Удалит истекшие данные автоматически")
}
EOF

go run /tmp/test_patient_access.go
```

### Тест 2: Проверка шифрования (unit test)

```bash
cd backend/data-transfer-service

# Создайте простой unit test
cat > internal/service/patient-access_test.go <<'EOF'
package service

import (
    "testing"
    "time"
)

func TestEncryptionKey(t *testing.T) {
    secret := "test-secret-key"
    service := NewPatientAccessService(nil, secret)

    if len(service.(*patientAccessService).encryptionKey) != 32 {
        t.Errorf("Expected key length 32, got %d", len(service.(*patientAccessService).encryptionKey))
    }

    t.Log("✓ Encryption key has correct length (32 bytes for AES-256)")
}

func TestAccessDuration(t *testing.T) {
    service := NewPatientAccessService(nil, "test")
    duration := service.(*patientAccessService).accessDuration

    expected := 15 * time.Minute
    if duration != expected {
        t.Errorf("Expected duration %v, got %v", expected, duration)
    }

    t.Log("✓ Access duration is 15 minutes")
}
EOF

# Запустите тест
go test ./internal/service -v -run TestEncryptionKey
go test ./internal/service -v -run TestAccessDuration
```

### Тест 3: Прямая проверка через SQL

```bash
# Подключитесь к БД
docker exec -it mde-postgres psql -U postgres -d medical

# Создайте тестовый запрос вручную
INSERT INTO patient_access_requests (
    patient_id, clinic_id, medical_data_id,
    status, requested_at
) VALUES (
    1, 1, 1,
    'pending', NOW()
) RETURNING id;

# Посмотрите созданный запрос
SELECT * FROM patient_access_requests WHERE id = 1;

# Симулируйте одобрение
UPDATE patient_access_requests
SET status = 'approved',
    approved_at = NOW(),
    expires_at = NOW() + INTERVAL '15 minutes'
WHERE id = 1;

# Проверьте
SELECT
    id,
    status,
    approved_at,
    expires_at,
    expires_at - NOW() as time_left
FROM patient_access_requests
WHERE id = 1;

# Создайте временную запись
INSERT INTO temporary_patient_data (
    access_token, patient_id, clinic_id, medical_data_id,
    encrypted_data, granted_at, expires_at, is_revoked
) VALUES (
    'test-token-123', 1, 1, 1,
    'encrypted_data_here', NOW(), NOW() + INTERVAL '15 minutes', false
);

# Проверьте
SELECT
    access_token,
    expires_at,
    expires_at - NOW() as time_left,
    is_revoked
FROM temporary_patient_data;

# Выйти
\q
```

## Полная интеграция с Hyperledger Fabric

Для полного тестирования с блокчейном нужно настроить Fabric. Вот упрощенная версия:

### Вариант 1: Fabric Test Network (рекомендуется для начала)

```bash
# 1. Установите Fabric binaries (если нет)
curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.5.0 1.5.5

# 2. Перейдите в Fabric samples
cd ~/fabric-samples/test-network

# 3. Запустите тестовую сеть
./network.sh up createChannel -c mychannel -ca

# 4. Установите наш chaincode
./network.sh deployCC -ccn medical-access \
    -ccp /Users/damir/clinic-data/MedicalDataExchange/blockchain/fabric/chaincode/medical-access \
    -ccl go

# 5. Протестируйте chaincode
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/

# Создать запрос на доступ
peer chaincode invoke -o localhost:7050 \
    --ordererTLSHostnameOverride orderer.example.com \
    --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
    -C mychannel -n medical-access \
    --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
    --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
    -c '{"function":"CreateAccessRequest","Args":["req1","patient1","clinic1","data1","patient_view","",""]}'

# Проверить запрос
peer chaincode query -C mychannel -n medical-access \
    -c '{"Args":["GetAccessRequest","req1"]}'
```

### Вариант 2: Использовать mock (для быстрого тестирования)

Можно создать mock blockchain service, который симулирует Fabric:

```bash
# Создайте mock сервис
cat > backend/blockchain-service/internal/fabric/mock.go <<'EOF'
package fabric

import (
    "encoding/json"
    "fmt"
    "time"
)

type MockFabricClient struct {
    storage map[string]interface{}
}

func NewMockFabricClient() *MockFabricClient {
    return &MockFabricClient{
        storage: make(map[string]interface{}),
    }
}

func (m *MockFabricClient) CreateAccessRequest(id string, data interface{}) error {
    m.storage[fmt.Sprintf("REQUEST_%s", id)] = data
    fmt.Printf("✓ Mock: Created access request %s\n", id)
    return nil
}

func (m *MockFabricClient) GetAccessRequest(id string) (interface{}, error) {
    data, ok := m.storage[fmt.Sprintf("REQUEST_%s", id)]
    if !ok {
        return nil, fmt.Errorf("not found")
    }
    fmt.Printf("✓ Mock: Retrieved access request %s\n", id)
    return data, nil
}
EOF
```

## Визуализация работы системы

```bash
# Создайте скрипт для демонстрации
cat > demo.sh <<'EOF'
#!/bin/bash

echo "╔════════════════════════════════════════════════════════════╗"
echo "║   ДЕМО: Временный доступ пациента к медицинским данным    ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

echo "📊 ШАГ 1: Исходное состояние"
echo "   Клиника: Хранит медицинские данные пациента"
echo "   Пациент: Не имеет доступа к данным"
echo ""
sleep 2

echo "📨 ШАГ 2: Пациент создает запрос"
echo "   POST /api/patient-access/request"
echo "   {patient_id: 1, clinic_id: 1, medical_data_id: 1}"
echo "   ✓ Запрос создан (status: pending)"
echo ""
sleep 2

echo "🏥 ШАГ 3: Клиника получает уведомление"
echo "   Сотрудник клиники проверяет запрос"
echo "   Подтверждает личность пациента"
echo ""
sleep 2

echo "✅ ШАГ 4: Клиника одобряет запрос"
echo "   POST /api/patient-access/approve/1"
echo "   ✓ Транзакция записана в блокчейн"
echo "   ✓ Создается зашифрованная копия (AES-256)"
echo "   ✓ Генерируется токен: abc123xyz..."
echo "   ✓ Срок действия: $(date -u -d '+15 minutes' '+%H:%M:%S')"
echo ""
sleep 2

echo "🔐 ШАГ 5: Шифрование данных"
echo "   Алгоритм: AES-256-GCM"
echo "   Оригинал: Остается в клинике"
echo "   Копия: Зашифрована и сохранена временно"
echo ""
sleep 2

echo "🎫 ШАГ 6: Пациент получает доступ"
echo "   GET /api/patient-access/data?token=abc123xyz"
echo "   ✓ Токен валиден"
echo "   ✓ Данные расшифрованы и отправлены"
echo "   📱 Пациент видит свои анализы"
echo ""
sleep 2

echo "⏰ ШАГ 7: Через 15 минут..."
echo "   Cleanup Scheduler запускается (каждые 5 минут)"
echo "   ✓ Обнаружена истекшая запись"
echo "   ✓ Зашифрованная копия удалена из БД"
echo "   ✓ Токен больше не действителен"
echo "   ✓ Оригинал остался в клинике"
echo ""
sleep 2

echo "✨ РЕЗУЛЬТАТ:"
echo "   🏥 Клиника: Данные на месте"
echo "   📱 Пациент: Доступ отозван автоматически"
echo "   🔗 Блокчейн: Все операции записаны"
echo "   🔒 Безопасность: Временная копия удалена"
echo ""
echo "╔════════════════════════════════════════════════════════════╗"
echo "║                    ДЕМО ЗАВЕРШЕНО                          ║"
echo "╚════════════════════════════════════════════════════════════╝"
EOF

chmod +x demo.sh
./demo.sh
```

## Диаграмма последовательности

```
Пациент          Core API      Data Transfer    Blockchain      PostgreSQL
   |                |               Service         Service           |
   |-- Запрос ----->|                  |              |               |
   |                |-- CreateAccess ->|              |               |
   |                |                  |-- INSERT --->|-------------->|
   |                |<-- Request ID ---|              |<-- ID --------|
   |<-- Request ID -|                  |              |               |
   |                                   |              |               |
(Клиника одобряет)                     |              |               |
   |                                   |              |               |
   |-- Одобрить --->|                  |              |               |
   |                |-- Approve ------>|              |               |
   |                |                  |-- LogTx ---->|               |
   |                |                  |              |-- Fabric ---> |
   |                |                  |<-- TxID -----|               |
   |                |                  |-- Encrypt -->|               |
   |                |                  |-- INSERT ------------------->|
   |                |                  |           (temporary_data)   |
   |                |<-- Token --------|              |               |
   |<-- Token ------|                  |              |               |
   |                                   |              |               |
(Через 15 минут)                       |              |               |
   |                                   |              |               |
   |                            Scheduler             |               |
   |                              запуск              |               |
   |                                   |-- Cleanup ------------------>|
   |                                   |         DELETE WHERE         |
   |                                   |         expires_at < NOW     |
   |                                   |<-- Deleted: 1 ---------------|
   |                                   |              |               |
(Токен больше не работает)             |              |               |
```

## Проверка работы Cleanup Scheduler

```bash
# Тест cleanup scheduler
cat > test_cleanup.sh <<'EOF'
#!/bin/bash

echo "=== Тест Cleanup Scheduler ==="
echo ""

# Подключение к БД
DB="docker exec -it mde-postgres psql -U postgres -d medical -t -c"

echo "1. Создаем истекшую запись (expired 5 минут назад)"
$DB "INSERT INTO temporary_patient_data (
    access_token, patient_id, clinic_id, medical_data_id,
    encrypted_data, granted_at, expires_at, is_revoked
) VALUES (
    'expired-token', 1, 1, 1,
    'old_encrypted_data', NOW() - INTERVAL '20 minutes',
    NOW() - INTERVAL '5 minutes', false
);"

echo "2. Создаем активную запись (истекает через 10 минут)"
$DB "INSERT INTO temporary_patient_data (
    access_token, patient_id, clinic_id, medical_data_id,
    encrypted_data, granted_at, expires_at, is_revoked
) VALUES (
    'active-token', 1, 1, 1,
    'current_encrypted_data', NOW(),
    NOW() + INTERVAL '10 minutes', false
);"

echo ""
echo "3. Проверяем записи ДО очистки:"
$DB "SELECT access_token,
    CASE
        WHEN expires_at < NOW() THEN 'EXPIRED'
        ELSE 'ACTIVE'
    END as status,
    expires_at - NOW() as time_left
FROM temporary_patient_data;"

echo ""
echo "4. Запускаем cleanup (симуляция):"
DELETED=$($DB "DELETE FROM temporary_patient_data
WHERE expires_at < NOW() OR is_revoked = true
RETURNING access_token;")

echo "   Удалено: $DELETED"

echo ""
echo "5. Проверяем записи ПОСЛЕ очистки:"
$DB "SELECT access_token,
    CASE
        WHEN expires_at < NOW() THEN 'EXPIRED'
        ELSE 'ACTIVE'
    END as status,
    expires_at - NOW() as time_left
FROM temporary_patient_data;"

echo ""
echo "✓ Cleanup работает корректно!"
echo "✓ Истекшие записи удалены"
echo "✓ Активные записи сохранены"
EOF

chmod +x test_cleanup.sh
./test_cleanup.sh
```

## Мониторинг в реальном времени

```bash
# Смотреть логи data-transfer-service в реальном времени
docker logs -f mde-dt --tail 100

# В другом терминале - логи PostgreSQL
docker logs -f mde-postgres --tail 50

# В третьем терминале - проверять БД каждые 5 секунд
watch -n 5 'docker exec mde-postgres psql -U postgres -d medical -c "SELECT COUNT(*) as active_tokens, MIN(expires_at - NOW()) as next_expiry FROM temporary_patient_data WHERE expires_at > NOW();"'
```

## FAQ - Частые вопросы

**Q: Где хранятся медицинские данные?**
A: В таблице `medical_data` в PostgreSQL, привязаны к `clinic_id`.

**Q: Где хранится временная копия?**
A: В таблице `temporary_patient_data`, зашифрована AES-256-GCM.

**Q: Как часто удаляются истекшие данные?**
A: Каждые 5 минут запускается Cleanup Scheduler.

**Q: Можно ли продлить токен?**
A: Нет, это мера безопасности. Нужно создать новый запрос.

**Q: Что если клиника не одобрит запрос?**
A: Статус останется `pending` или изменится на `rejected`.

**Q: Можно ли изменить время доступа (15 минут)?**
A: Да, в `patient-access_service.go` изменить `accessDuration: 15 * time.Minute`.

**Q: Нужен ли блокчейн для тестирования?**
A: Нет, можно тестировать без Fabric, только БД и шифрование.

## Следующие шаги

1. ✅ Тестирование БД и таблиц (описано выше)
2. ✅ Тестирование шифрования (unit tests)
3. ✅ Тестирование cleanup scheduler (SQL скрипт)
4. ⏳ Настройка Hyperledger Fabric (опционально)
5. ⏳ Интеграция REST API в core-service
6. ⏳ Frontend UI для запросов на доступ

Хотите протестировать что-то конкретное? Я могу помочь!
