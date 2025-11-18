как# Блокчейн-система управления доступом к медицинским данным

## Обзор архитектуры

Система использует **Hyperledger Fabric** для управления доступом к медицинским данным с временными правами доступа.

### Ключевые компоненты:

1. **Hyperledger Fabric** - консорциумный блокчейн для аудита и управления правами
2. **Chaincode (Smart Contract)** - логика управления доступом на Go
3. **Blockchain Service** - Go-сервис для взаимодействия с Fabric
4. **Data Transfer Service** - управление временными данными и шифрование
5. **PostgreSQL** - хранение зашифрованных временных копий данных

## Основные сценарии использования

### 1. Пациент запрашивает доступ к своим данным

**Workflow:**

```
1. Пациент → Создает запрос на просмотр медицинских данных
2. Запрос сохраняется в БД (status: pending)
3. Клиника → Получает уведомление о запросе
4. Клиника → Одобряет запрос
5. Система:
   - Записывает транзакцию в блокчейн
   - Создает зашифрованную копию данных (AES-256-GCM)
   - Генерирует временный токен доступа
   - Устанавливает срок действия: 15 минут
6. Пациент → Получает доступ к данным по токену
7. Через 15 минут:
   - Токен истекает
   - Зашифрованная копия удаляется автоматически
   - Данные остаются в клинике
```

**API Endpoints:**

```bash
# Создать запрос на доступ (пациент)
POST /api/patient-access/request
{
  "patient_id": 123,
  "clinic_id": 1,
  "medical_data_id": 456
}

# Одобрить запрос (клиника)
POST /api/patient-access/approve/{request_id}
{
  "approver_id": 789  # ID сотрудника клиники
}

# Получить данные по токену (пациент)
GET /api/patient-access/data?token=<access_token>
# Доступно только в течение 15 минут!
```

### 2. Передача данных между клиниками

**Workflow:**

```
1. Клиника A → Создает запрос на передачу данных в Клинику B
2. Пациент → Получает уведомление
3. Пациент → Одобряет передачу
4. Система:
   - Записывает транзакцию в блокчейн с hash данных
   - Копирует данные в Клинику B
   - Данные остаются в обеих клиниках
5. Блокчейн → Хранит audit trail передачи
```

**API Endpoints:**

```bash
# Создать запрос на передачу (Клиника A)
POST /api/transfers
{
  "medical_data_id": 456,
  "from_clinic_id": 1,
  "to_clinic_id": 2,
  "user_id": 123
}

# Одобрить передачу (пациент)
POST /api/transfers/decision
{
  "transfer_id": 789,
  "user_id": 123,
  "confirm": true
}
```

## Блокчейн Chaincode (Smart Contract)

### Основные функции:

#### 1. CreateAccessRequest
Создает запрос на временный доступ к данным.

```go
CreateAccessRequest(
    id string,
    patientID string,
    clinicID string,
    medicalDataID string,
    requestType string,  // patient_view или clinic_transfer
    fromClinicID string,
    toClinicID string
)
```

#### 2. ApproveAccessRequest
Одобряет запрос и создает временный токен (15 минут).

```go
ApproveAccessRequest(
    requestID string,
    approverID string
)
```

#### 3. ValidateAccess
Проверяет действительность токена доступа.

```go
ValidateAccess(accessID string) -> TemporaryAccess
```

#### 4. RevokeAccess
Досрочно отзывает доступ.

```go
RevokeAccess(accessID string)
```

#### 5. LogDataTransfer
Записывает передачу данных между клиниками.

```go
LogDataTransfer(
    id string,
    medicalDataID string,
    patientID string,
    fromClinicID string,
    toClinicID string,
    approvedBy string,
    dataHash string  // SHA-256 hash для аудита
)
```

## Шифрование данных

### Алгоритм: AES-256-GCM

**Процесс шифрования временных данных:**

1. Медицинские данные сериализуются в JSON
2. Генерируется уникальный nonce (12 bytes)
3. Данные шифруются с помощью AES-256-GCM
4. Результат кодируется в Base64
5. Сохраняется в таблице `temporary_patient_data`

**Ключ шифрования:**
- Генерируется из секрета через SHA-256
- Хранится в переменных окружения
- Используется только на стороне сервера

```go
// Пример генерации ключа
hash := sha256.Sum256([]byte(encryptionSecret))
encryptionKey := hash[:]  // 32 bytes для AES-256
```

## База данных

### Таблица: patient_access_requests

```sql
CREATE TABLE patient_access_requests (
    id BIGSERIAL PRIMARY KEY,
    patient_id BIGINT NOT NULL,
    clinic_id BIGINT NOT NULL,
    medical_data_id BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL, -- pending, approved, rejected, expired
    blockchain_tx_id VARCHAR(255),
    requested_at TIMESTAMP NOT NULL,
    approved_at TIMESTAMP,
    expires_at TIMESTAMP
);
```

### Таблица: temporary_patient_data

```sql
CREATE TABLE temporary_patient_data (
    id BIGSERIAL PRIMARY KEY,
    access_token VARCHAR(255) UNIQUE NOT NULL,
    patient_id BIGINT NOT NULL,
    clinic_id BIGINT NOT NULL,
    medical_data_id BIGINT NOT NULL,
    encrypted_data TEXT NOT NULL,  -- Зашифрованные данные
    granted_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,  -- +15 минут от granted_at
    is_revoked BOOLEAN DEFAULT FALSE
);
```

## Автоматическая очистка данных

### Cleanup Scheduler

Запускается каждые **5 минут** и:

1. Удаляет все записи из `temporary_patient_data` где:
   - `expires_at < NOW()` ИЛИ
   - `is_revoked = true`

2. Обновляет статус в `patient_access_requests`:
   - `status = 'expired'` для истекших запросов

```go
// Запуск scheduler при старте сервиса
scheduler := scheduler.NewCleanupScheduler(patientAccessService)
go scheduler.Start(context.Background())
```

## Безопасность

### 1. Шифрование
- **AES-256-GCM** для временных данных
- Уникальный nonce для каждого шифрования
- Authenticated encryption (защита от модификации)

### 2. Временные токены
- Генерируются с помощью SHA-256
- Содержат timestamp для уникальности
- Автоматически истекают через 15 минут
- Невозможно подделать без ключа

### 3. Блокчейн аудит
- Все запросы и передачи записываются в Fabric
- Immutable audit trail
- Hash данных для верификации целостности
- Консенсус между клиниками

### 4. Контроль доступа
- Пациент видит только свои данные
- Клиника одобряет запросы только для своих пациентов
- RBAC на уровне gRPC interceptors

## Запуск системы

### 1. Запуск Fabric Network

```bash
cd blockchain/fabric/network

# Генерация криптоматериалов (первый раз)
./generate-crypto.sh

# Запуск сети
docker-compose -f docker-compose-fabric.yaml up -d

# Создание канала и установка chaincode
./setup-network.sh
```

### 2. Запуск бэкенда

```bash
# Из корневой директории
docker-compose up -d

# Сервисы:
# - postgres (5442)
# - user-service (50053)
# - data-transfer-service (50052)
# - blockchain-service (50054)
# - core-service (8099)
# - frontend (8080)
```

### 3. Проверка работы

```bash
# Проверить статус Fabric
docker exec -it fabric-cli peer channel list

# Проверить chaincode
docker exec -it fabric-cli peer chaincode query \
  -C medical-channel \
  -n medical-access \
  -c '{"Args":["GetAccessRequest","request_1"]}'
```

## Мониторинг

### Логи блокчейн-транзакций

```bash
# Логи blockchain-service
docker logs -f mde-blockchain

# Логи peer
docker logs -f peer0.clinic1.medical.com

# Логи orderer
docker logs -f orderer.medical.com
```

### Метрики

```bash
# Prometheus endpoint
curl http://localhost:9444/metrics  # Peer0 Clinic1
curl http://localhost:9445/metrics  # Peer0 Clinic2
```

## Примеры использования

### Пример 1: Пациент просматривает свои данные

```javascript
// 1. Пациент создает запрос
const response = await fetch('http://localhost:8099/api/patient-access/request', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer ' + patientToken,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    patient_id: 123,
    clinic_id: 1,
    medical_data_id: 456
  })
});

const { request_id } = await response.json();

// 2. Клиника одобряет запрос
const approveResp = await fetch(`http://localhost:8099/api/patient-access/approve/${request_id}`, {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer ' + clinicToken,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    approver_id: 789
  })
});

const { access_token, expires_at } = await approveResp.json();

// 3. Пациент получает данные (в течение 15 минут!)
const dataResp = await fetch(`http://localhost:8099/api/patient-access/data?token=${access_token}`, {
  headers: {
    'Authorization': 'Bearer ' + patientToken
  }
});

const medicalData = await dataResp.json();
console.log('Данные доступны до:', expires_at);
```

### Пример 2: Передача данных между клиниками

```javascript
// 1. Клиника A создает запрос на передачу
const transferReq = await fetch('http://localhost:8099/api/transfers', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer ' + clinic1Token,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    medical_data_id: 456,
    from_clinic_id: 1,
    to_clinic_id: 2,
    user_id: 123
  })
});

const { transfer_id } = await transferReq.json();

// 2. Пациент одобряет передачу
const decisionResp = await fetch('http://localhost:8099/api/transfers/decision', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer ' + patientToken,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    transfer_id: transfer_id,
    user_id: 123,
    confirm: true
  })
});

// 3. Данные теперь доступны в обеих клиниках
// Транзакция записана в блокчейн для аудита
```

## Troubleshooting

### Проблема: Токен доступа истек

```bash
# Решение: Создать новый запрос на доступ
# Токены нельзя продлить - это мера безопасности
```

### Проблема: Ошибка при расшифровке данных

```bash
# Проверить переменную окружения ENCRYPTION_SECRET
# Она должна быть одинаковой при шифровании и расшифровке
```

### Проблема: Chaincode не отвечает

```bash
# Проверить статус peers
docker ps | grep peer

# Переустановить chaincode
docker exec -it fabric-cli peer chaincode install ...
docker exec -it fabric-cli peer chaincode approve ...
docker exec -it fabric-cli peer chaincode commit ...
```

## Дополнительные возможности

### Расширение срока доступа

Для изменения времени доступа (по умолчанию 15 минут):

```go
// В blockchain/fabric/chaincode/medical-access/main.go
expiresAt := now.Add(30 * time.Minute) // Изменить на 30 минут

// В backend/data-transfer-service/internal/service/patient-access_service.go
accessDuration: 30 * time.Minute // Изменить на 30 минут
```

### Добавление аудита действий

Все действия можно отслеживать через блокчейн:

```bash
# Получить историю доступа пациента
docker exec -it fabric-cli peer chaincode query \
  -C medical-channel \
  -n medical-access \
  -c '{"Args":["GetPatientAccessHistory","patient_123"]}'
```

## Заключение

Эта система обеспечивает:
- ✅ Временный доступ пациентов к их данным (10-15 минут)
- ✅ Автоматическое удаление временных данных
- ✅ Данные хранятся у клиники, не у пациента
- ✅ Безопасную передачу между клиниками с одобрением пациента
- ✅ Полный audit trail в блокчейне
- ✅ Шифрование временных копий (AES-256)
- ✅ Невозможность подделки или продления токенов

**Контакты:** Для вопросов и предложений создавайте issues в репозитории.
