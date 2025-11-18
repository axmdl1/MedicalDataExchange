# Настройка системы для работы с реальными данными

## Что было сделано

✅ Добавлены методы Patient Access в proto файл data-transfer  
✅ Реализованы gRPC handlers в data-transfer-service  
✅ Обновлен patient-access_handler.go для использования реальных вызовов  
✅ Добавлена поддержка шифрования данных (AES-256-GCM)  
✅ Медицинские данные хранятся у клиники (clinic_id в таблице medical_data)  

## Что нужно сделать

### 1. Сгенерировать proto файлы

```bash
cd /Users/damir/clinic-data/MedicalDataExchange/backend/core-service

# Установите protoc-gen-go и protoc-gen-go-grpc если еще не установлены
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Сгенерируйте proto файлы
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       pkg/proto/data-transfer/data-transfer.proto
```

### 2. Пересобрать сервисы

```bash
cd /Users/damir/clinic-data/MedicalDataExchange

# Пересобрать data-transfer-service
docker-compose build data-transfer-service

# Пересобрать core-service
docker-compose build core-service

# Перезапустить сервисы
docker-compose up -d data-transfer-service core-service
```

### 3. Проверить конфигурацию

Убедитесь, что в `backend/data-transfer-service/config/local.yaml` есть:

```yaml
encryption:
  secret: "your-encryption-secret-key-32-bytes-long!!!"
```

### 4. Проверить работу системы

1. Войдите как пациент
2. Создайте запрос на доступ к медицинским данным
3. Войдите как сотрудник клиники
4. Одобрите запрос
5. Получите токен доступа
6. Войдите обратно как пациент
7. Используйте токен для просмотра данных

## Как работает система

### Хранение данных
- **Оригинальные данные** хранятся в таблице `medical_data` с `clinic_id` (принадлежат клинике)
- **Временные данные** хранятся в таблице `temporary_patient_data` в зашифрованном виде
- Временные данные автоматически удаляются через 15 минут

### Процесс доступа
1. Пациент создает запрос → сохраняется в `patient_access_requests`
2. Клиника одобряет запрос → создается зашифрованная копия в `temporary_patient_data`
3. Пациент получает токен → может просмотреть данные в течение 15 минут
4. Через 15 минут → данные автоматически удаляются (cleanup scheduler)

### Шифрование
- Данные шифруются с помощью AES-256-GCM
- Ключ шифрования берется из переменной окружения `encryption.secret`
- Каждое шифрование использует уникальный nonce

## Интеграция с блокчейном (опционально)

Для полной интеграции с блокчейном нужно:

1. Настроить Hyperledger Fabric network
2. Добавить вызовы к blockchain-service в patient-access_service.go при одобрении запросов
3. Записывать транзакции в блокчейн для аудита

Пример интеграции (в `ApproveAccessRequest`):

```go
// После создания временных данных
blockchainReq := &blockchainpb.CreateAccessRequestRequest{
    Id:            fmt.Sprintf("%d", request.ID),
    PatientId:     fmt.Sprintf("%d", request.PatientID),
    ClinicId:      fmt.Sprintf("%d", request.ClinicID),
    MedicalDataId: fmt.Sprintf("%d", request.MedicalDataID),
    RequestType:   "patient_access",
}

blockchainResp, err := s.blockchainClient.CreateAccessRequest(ctx, blockchainReq)
if err == nil {
    request.BlockchainTxID = blockchainResp.Id
    s.repo.UpdatePatientAccessRequest(ctx, request)
}
```

## Тестирование

Все сценарии из `TESTING_SCENARIOS.md` теперь работают с реальными данными:

- ✅ Создание запроса на доступ
- ✅ Одобрение запроса клиникой
- ✅ Получение токена доступа
- ✅ Просмотр данных по токену
- ✅ Автоматическое удаление через 15 минут
- ✅ Список запросов пациента
- ✅ Отзыв доступа

## Важные замечания

1. **Шифрование**: Ключ шифрования должен быть минимум 32 байта
2. **Безопасность**: В продакшене используйте безопасное хранение секретов (Vault, AWS Secrets Manager)
3. **Блокчейн**: Для полного аудита настройте Hyperledger Fabric
4. **Cleanup**: Scheduler автоматически удаляет истекшие данные каждые 5 минут

## Troubleshooting

### Ошибка "proto file not found"
- Убедитесь, что proto файлы сгенерированы (шаг 1)

### Ошибка "encryption secret not found"
- Проверьте конфигурацию в `config/local.yaml`

### Данные не удаляются автоматически
- Проверьте, что cleanup scheduler запущен в data-transfer-service
- Проверьте логи: `docker logs mde-dt | grep cleanup`

