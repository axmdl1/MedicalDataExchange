# 🔧 Fixing Mock Data Issue

## Problem

В dropdown показывается "test (13.11.2025)" вместо реальных диагнозов типа "Acute URTI (J06.9)".

**Причина:** В БД остались старые мок данные.

## Quick Fix (Рекомендуется)

### Вариант 1: Автоматический cleanup (БЕЗ остановки системы)

```bash
# Запустите скрипт
./cleanup-and-reload.sh
```

**Что делает скрипт:**
1. Удаляет записи с `diagnosis = 'test'`
2. Удаляет записи без visit_date и department
3. Перезагружает реальные медицинские данные
4. Показывает текущие записи

**После запуска обновите браузер!**

### Вариант 2: Полный пересоздание БД

```bash
# Остановить и удалить данные
docker-compose down -v

# Запустить заново (создаст БД с нуля)
docker-compose up -d

# Подождать 10-15 секунд для инициализации
sleep 15

# Проверить данные
docker exec -i mde-postgres psql -U postgres -d medical -c \
  "SELECT id, diagnosis, diagnosis_code FROM medical_data;"
```

## Verification

После исправления в dropdown должны показываться:

✅ **Правильно:**
```
Acute upper respiratory tract infection [J06.9] - General Medicine (16.11.2025)
Chronic gastritis with reflux esophagitis [K29.5] - Gastroenterology (14.11.2025)
Essential (primary) hypertension [I10] - Cardiology (30.10.2025)
Allergic rhinitis, seasonal [J30.1] - Allergy & Immunology (12.11.2025)
```

❌ **Неправильно:**
```
test (13.11.2025)
```

## Check Current Data

Проверить текущие данные в БД:

```bash
docker exec -i mde-postgres psql -U postgres -d medical <<EOF
SELECT
    id,
    user_id,
    diagnosis,
    diagnosis_code,
    department,
    visit_date::date,
    created_at::date
FROM medical_data
ORDER BY visit_date DESC;
EOF
```

## Expected Output

```
 id | user_id | diagnosis                                  | diagnosis_code | department           | visit_date | created_at
----+---------+--------------------------------------------+----------------+----------------------+------------+------------
  1 |       1 | Acute upper respiratory tract infection    | J06.9          | General Medicine     | 2025-11-16 | 2025-11-19
  2 |       1 | Healthy individual - routine checkup       | Z00.0          | General Medicine     | 2025-10-15 | 2025-11-19
  3 |       2 | Chronic gastritis with reflux esophagitis  | K29.5          | Gastroenterology     | 2025-11-14 | 2025-11-19
  4 |       2 | Essential (primary) hypertension, Stage 2  | I10            | Cardiology           | 2025-10-30 | 2025-11-19
  5 |       3 | Allergic rhinitis, seasonal (hay fever)    | J30.1          | Allergy & Immunology | 2025-11-12 | 2025-11-19
  6 |       3 | Sprain of lateral ligament of right ankle  | S93.41         | Orthopedics          | 2025-11-07 | 2025-11-19
```

## Manual Cleanup (если скрипт не работает)

```bash
# Подключиться к БД
docker exec -it mde-postgres psql -U postgres -d medical

# Удалить мок данные
DELETE FROM medical_data WHERE diagnosis = 'test';
DELETE FROM medical_data WHERE visit_date IS NULL AND department IS NULL;

# Проверить что осталось
SELECT COUNT(*) FROM medical_data;

# Если пусто - загрузить реальные данные
\i /docker-entrypoint-initdb.d/21_real_medical_data.sql

# Выход
\q
```

## Frontend Update

Frontend уже обновлен! Теперь показывает:
```javascript
`${diagnosis} [${diagnosis_code}] - ${department} (${date})`
```

Примеры:
- "Acute upper respiratory tract infection [J06.9] - General Medicine (16.11.2025)"
- "Chronic gastritis with reflux esophagitis [K29.5] - Gastroenterology (14.11.2025)"

## Troubleshooting

**Проблема:** После cleanup все еще показывает "test"

**Решение:**
1. Hard refresh браузера: Ctrl+Shift+R (Windows) / Cmd+Shift+R (Mac)
2. Очистить кэш браузера
3. Проверить что данные обновились в БД (см. Check Current Data выше)

**Проблема:** Dropdown пустой

**Решение:**
1. Проверить что пациент имеет записи в выбранной клинике
2. Проверить patient_id и clinic_id совпадают:
   ```sql
   SELECT * FROM medical_data WHERE user_id = 1 AND clinic_id = 1;
   ```

## Summary

- ✅ Создан `cleanup-and-reload.sh` - автоматическая очистка
- ✅ Создан `22_cleanup_mock_data.sql` - SQL для очистки
- ✅ Frontend обновлен - показывает diagnosis + код + отделение + дата
- ✅ После fix dropdown покажет реальные диагнозы

**Запустите:** `./cleanup-and-reload.sh` и обновите браузер!
