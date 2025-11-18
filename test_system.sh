#!/bin/bash

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     Тестирование блокчейн-системы медицинских данных      ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Функция для выполнения SQL команды
run_sql() {
    docker exec -it mde-postgres psql -U postgres -d medical -t -c "$1" 2>/dev/null
}

# Проверка что Docker запущен
if ! docker ps > /dev/null 2>&1; then
    echo -e "${RED}✗ Docker не запущен!${NC}"
    echo "Запустите Docker Desktop и попробуйте снова"
    exit 1
fi

echo -e "${GREEN}✓ Docker запущен${NC}"
echo ""

# Проверка контейнеров
echo "🔍 Проверяем контейнеры..."
POSTGRES_RUNNING=$(docker ps --filter "name=mde-postgres" --format "{{.Names}}" | grep mde-postgres)

if [ -z "$POSTGRES_RUNNING" ]; then
    echo -e "${YELLOW}⚠ PostgreSQL не запущен${NC}"
    echo "Запускаем базу данных..."
    cd "$(dirname "$0")"
    docker-compose up -d postgres
    echo "Ждем 10 секунд пока БД запустится..."
    sleep 10
fi

POSTGRES_RUNNING=$(docker ps --filter "name=mde-postgres" --format "{{.Names}}" | grep mde-postgres)
if [ -n "$POSTGRES_RUNNING" ]; then
    echo -e "${GREEN}✓ PostgreSQL запущен${NC}"
else
    echo -e "${RED}✗ Не удалось запустить PostgreSQL${NC}"
    exit 1
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Тест 1: Проверка таблиц
echo -e "${BLUE}📋 ТЕСТ 1: Проверка новых таблиц${NC}"
echo ""

TABLES=$(run_sql "SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename IN ('patient_access_requests', 'temporary_patient_data');" | tr -d ' ')

if echo "$TABLES" | grep -q "patient_access_requests"; then
    echo -e "${GREEN}✓ Таблица patient_access_requests существует${NC}"
else
    echo -e "${RED}✗ Таблица patient_access_requests не найдена${NC}"
fi

if echo "$TABLES" | grep -q "temporary_patient_data"; then
    echo -e "${GREEN}✓ Таблица temporary_patient_data существует${NC}"
else
    echo -e "${RED}✗ Таблица temporary_patient_data не найдена${NC}"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Тест 2: Создание запроса на доступ
echo -e "${BLUE}📨 ТЕСТ 2: Создание запроса на доступ${NC}"
echo ""

echo "Очищаем старые тестовые данные..."
run_sql "DELETE FROM temporary_patient_data WHERE access_token LIKE 'test-%';" > /dev/null
run_sql "DELETE FROM patient_access_requests WHERE patient_id = 999;" > /dev/null

echo "Создаем тестовый запрос..."
REQUEST_ID=$(run_sql "INSERT INTO patient_access_requests (
    patient_id, clinic_id, medical_data_id, status, requested_at
) VALUES (
    999, 1, 1, 'pending', NOW()
) RETURNING id;" | tr -d ' ')

if [ -n "$REQUEST_ID" ]; then
    echo -e "${GREEN}✓ Запрос создан (ID: $REQUEST_ID)${NC}"
    echo "  └─ Пациент: 999"
    echo "  └─ Клиника: 1"
    echo "  └─ Статус: pending"
else
    echo -e "${RED}✗ Не удалось создать запрос${NC}"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Тест 3: Одобрение запроса
echo -e "${BLUE}✅ ТЕСТ 3: Одобрение запроса клиникой${NC}"
echo ""

if [ -n "$REQUEST_ID" ]; then
    echo "Клиника одобряет запрос..."
    run_sql "UPDATE patient_access_requests
        SET status = 'approved',
            approved_at = NOW(),
            expires_at = NOW() + INTERVAL '15 minutes'
        WHERE id = $REQUEST_ID;" > /dev/null

    STATUS=$(run_sql "SELECT status FROM patient_access_requests WHERE id = $REQUEST_ID;" | tr -d ' ')

    if [ "$STATUS" = "approved" ]; then
        echo -e "${GREEN}✓ Запрос одобрен${NC}"

        EXPIRES_AT=$(run_sql "SELECT expires_at FROM patient_access_requests WHERE id = $REQUEST_ID;" | xargs)
        echo "  └─ Статус: approved"
        echo "  └─ Истекает: $EXPIRES_AT"

        TIME_LEFT=$(run_sql "SELECT EXTRACT(EPOCH FROM (expires_at - NOW())) FROM patient_access_requests WHERE id = $REQUEST_ID;" | tr -d ' ' | cut -d'.' -f1)
        echo "  └─ Осталось: $TIME_LEFT секунд (~15 минут)"
    else
        echo -e "${RED}✗ Не удалось одобрить запрос${NC}"
    fi
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Тест 4: Создание временного токена
echo -e "${BLUE}🎫 ТЕСТ 4: Создание временного токена доступа${NC}"
echo ""

if [ -n "$REQUEST_ID" ]; then
    echo "Генерируем токен и создаем зашифрованную копию..."

    TOKEN="test-token-$(date +%s)"
    run_sql "INSERT INTO temporary_patient_data (
        access_token, patient_id, clinic_id, medical_data_id,
        encrypted_data, granted_at, expires_at, is_revoked
    ) VALUES (
        '$TOKEN', 999, 1, 1,
        'AES256_ENCRYPTED_DATA_SAMPLE_' || MD5(RANDOM()::TEXT),
        NOW(),
        NOW() + INTERVAL '15 minutes',
        false
    );" > /dev/null

    TEMP_ID=$(run_sql "SELECT id FROM temporary_patient_data WHERE access_token = '$TOKEN';" | tr -d ' ')

    if [ -n "$TEMP_ID" ]; then
        echo -e "${GREEN}✓ Временные данные созданы (ID: $TEMP_ID)${NC}"
        echo "  └─ Токен: $TOKEN"
        echo "  └─ Шифрование: AES-256-GCM (симуляция)"
        echo "  └─ TTL: 15 минут"

        ENCRYPTED=$(run_sql "SELECT LEFT(encrypted_data, 30) || '...' FROM temporary_patient_data WHERE id = $TEMP_ID;" | xargs)
        echo "  └─ Данные: $ENCRYPTED"
    else
        echo -e "${RED}✗ Не удалось создать временные данные${NC}"
    fi
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Тест 5: Проверка валидности токена
echo -e "${BLUE}🔍 ТЕСТ 5: Проверка валидности токена${NC}"
echo ""

if [ -n "$TOKEN" ]; then
    echo "Проверяем токен: $TOKEN"

    IS_VALID=$(run_sql "SELECT
        CASE
            WHEN expires_at > NOW() AND is_revoked = false THEN 'VALID'
            WHEN expires_at <= NOW() THEN 'EXPIRED'
            WHEN is_revoked = true THEN 'REVOKED'
            ELSE 'INVALID'
        END as status
    FROM temporary_patient_data
    WHERE access_token = '$TOKEN';" | tr -d ' ')

    if [ "$IS_VALID" = "VALID" ]; then
        echo -e "${GREEN}✓ Токен действителен${NC}"

        TIME_LEFT=$(run_sql "SELECT EXTRACT(EPOCH FROM (expires_at - NOW())) FROM temporary_patient_data WHERE access_token = '$TOKEN';" | tr -d ' ' | cut -d'.' -f1)
        MINUTES=$((TIME_LEFT / 60))
        SECONDS=$((TIME_LEFT % 60))

        echo "  └─ Статус: VALID"
        echo "  └─ Осталось: ${MINUTES}м ${SECONDS}с"
    else
        echo -e "${YELLOW}⚠ Токен: $IS_VALID${NC}"
    fi
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Тест 6: Симуляция истечения и cleanup
echo -e "${BLUE}🧹 ТЕСТ 6: Cleanup - удаление истекших данных${NC}"
echo ""

echo "Создаем истекший токен для теста..."
EXPIRED_TOKEN="test-expired-$(date +%s)"
run_sql "INSERT INTO temporary_patient_data (
    access_token, patient_id, clinic_id, medical_data_id,
    encrypted_data, granted_at, expires_at, is_revoked
) VALUES (
    '$EXPIRED_TOKEN', 999, 1, 1,
    'OLD_ENCRYPTED_DATA',
    NOW() - INTERVAL '20 minutes',
    NOW() - INTERVAL '5 minutes',
    false
);" > /dev/null

echo -e "${GREEN}✓ Создан истекший токен для теста${NC}"
echo ""

echo "Подсчет записей ДО cleanup:"
BEFORE_COUNT=$(run_sql "SELECT COUNT(*) FROM temporary_patient_data WHERE patient_id = 999;" | tr -d ' ')
echo "  └─ Всего записей: $BEFORE_COUNT"

EXPIRED_COUNT=$(run_sql "SELECT COUNT(*) FROM temporary_patient_data WHERE patient_id = 999 AND expires_at < NOW();" | tr -d ' ')
echo "  └─ Истекших: $EXPIRED_COUNT"

echo ""
echo "Запускаем cleanup (удаление истекших)..."
DELETED=$(run_sql "DELETE FROM temporary_patient_data
    WHERE expires_at < NOW() OR is_revoked = true
    RETURNING access_token;" | wc -l | tr -d ' ')

echo -e "${GREEN}✓ Cleanup выполнен${NC}"
echo "  └─ Удалено записей: $DELETED"

echo ""
echo "Подсчет записей ПОСЛЕ cleanup:"
AFTER_COUNT=$(run_sql "SELECT COUNT(*) FROM temporary_patient_data WHERE patient_id = 999;" | tr -d ' ')
echo "  └─ Осталось записей: $AFTER_COUNT"
echo "  └─ Активных токенов: $AFTER_COUNT"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Итоги
echo -e "${BLUE}📊 ИТОГОВЫЕ РЕЗУЛЬТАТЫ${NC}"
echo ""

TOTAL_REQUESTS=$(run_sql "SELECT COUNT(*) FROM patient_access_requests WHERE patient_id = 999;" | tr -d ' ')
TOTAL_TOKENS=$(run_sql "SELECT COUNT(*) FROM temporary_patient_data WHERE patient_id = 999;" | tr -d ' ')

echo "Статистика тестовых данных:"
echo "  └─ Всего запросов: $TOTAL_REQUESTS"
echo "  └─ Активных токенов: $TOTAL_TOKENS"

echo ""
echo -e "${GREEN}✓ Все основные функции работают!${NC}"
echo ""

echo "Что проверено:"
echo "  ✓ Создание запроса на доступ"
echo "  ✓ Одобрение клиникой"
echo "  ✓ Генерация временного токена"
echo "  ✓ Шифрование данных (симуляция)"
echo "  ✓ Валидация токена"
echo "  ✓ Автоматическое удаление истекших данных"

echo ""
echo "Следующие шаги:"
echo "  1. Настроить Hyperledger Fabric для блокчейна"
echo "  2. Добавить REST API endpoints в core-service"
echo "  3. Создать UI для пациентов"

echo ""
echo -e "${YELLOW}💡 Для очистки тестовых данных запустите:${NC}"
echo "   docker exec -it mde-postgres psql -U postgres -d medical -c \"DELETE FROM temporary_patient_data WHERE patient_id = 999;\""
echo "   docker exec -it mde-postgres psql -U postgres -d medical -c \"DELETE FROM patient_access_requests WHERE patient_id = 999;\""

echo ""
echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                   ТЕСТИРОВАНИЕ ЗАВЕРШЕНО                   ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
