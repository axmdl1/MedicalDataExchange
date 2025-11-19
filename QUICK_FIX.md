# 🔧 Быстрое исправление Mock Data

## Проблема

В dropdown показывается **"test (13.11.2025)"** вместо реальных диагнозов.

## Решение (30 секунд)

```bash
# Запустите этот скрипт
./cleanup-and-reload.sh
```

**Что он делает:**
1. Удаляет записи с `diagnosis = 'test'`
2. Загружает 6 реальных медицинских записей
3. Показывает результат

## После исправления

**Обновите страницу в браузере (Ctrl+Shift+R)**

Теперь в dropdown должно показываться:

✅ **Пациент 1 (Anna Ivanova) в City Clinic:**
- Acute upper respiratory tract infection [J06.9] - General Medicine (16.11.2025)
- Healthy individual - routine checkup [Z00.0] - General Medicine (15.10.2025)

✅ **Пациент 2 (Boris Petrov) в Regional Hospital:**
- Chronic gastritis with reflux esophagitis [K29.5] - Gastroenterology (14.11.2025)
- Essential (primary) hypertension, Stage 2 [I10] - Cardiology (30.10.2025)

✅ **Пациент 3 (Svetlana Kuznetsova) в Medical Center Plus:**
- Allergic rhinitis, seasonal (hay fever) [J30.1] - Allergy & Immunology (12.11.2025)
- Sprain of lateral ligament of right ankle [S93.41] - Orthopedics (07.11.2025)

## Проверка

```bash
# Проверить данные в БД
docker exec -i mde-postgres psql -U postgres -d medical -c \
  "SELECT id, diagnosis, diagnosis_code, department FROM medical_data;"
```

Должны увидеть **6 реальных записей** с полными диагнозами, кодами МКБ-10 и отделениями.

## Если не помогло

**Полная перезагрузка БД:**
```bash
docker-compose down -v      # Удалит все данные
docker-compose up -d         # Создаст заново
sleep 15                     # Подождать инициализации
```

Реальные данные загрузятся автоматически из `21_real_medical_data.sql`.

## Готово!

После fix у вас будет полноценная система с реальными клиническими данными:
- ОРВИ с температурой 38.7°C, витальными показателями, анализами
- Гастрит с эндоскопией и H.pylori тестом
- Гипертония с кардиологическим обследованием
- Аллергический ринит с IgE тестами
- Спортивная травма с рентгеном

**Все данные реальные, не моки! 🎉**
