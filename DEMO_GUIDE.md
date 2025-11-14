# 🎬 Demo Guide - Полный сценарий работы системы

## Цель демонстрации
Показать полный цикл работы с медицинскими данными:
1. Врач создает пациента и его медицинские записи
2. Врач создает запрос на передачу данных в другую клинику
3. Пациент видит запрос и подтверждает передачу
4. Данные появляются в целевой клинике

---

## 📋 Подготовка

### 1. Запуск системы

```bash
# Запустить бэкенд
cd MedicalDataExchange
docker compose down -v  # Очистить старые данные
docker compose up -d    # Запустить с чистой БД

# Проверить что все работает
docker compose ps

# Запустить фронтенд (в другом терминале)
cd ..
python3 -m http.server 8000
```

Открыть: **http://localhost:8000**

### 2. Тестовые аккаунты

**Admin:**
- Email: `admin@system.com`
- Password: `admin123`

**Врач City Clinic (clinic_id = 1):**
- Email: `doctor1@clinic.com`
- Password: `doctor123`

**Врач Regional Hospital (clinic_id = 2):**
- Email: `doctor2@hospital.com`
- Password: `doctor123`

**Пациент:**
- Email: `patient1@example.com`
- Password: `patient123`

---

## 🎯 Сценарий демонстрации

### Шаг 1: Войти как врач (employee)

```
Email: doctor1@clinic.com
Password: doctor123
```

Вы видите:
- ✅ Dashboard с статистикой клиники
- ✅ Меню: Patients, Employees, Medical Records, Data Transfers
- ✅ Кнопки создания: Create User, Create Medical Data, Create Transfer

---

### Шаг 2: Создать нового пациента

1. Перейти в **Create User**
2. Заполнить форму:
   - User Type: `Patient`
   - First Name: `Тестовый`
   - Last Name: `Пациент`
   - Email: `testpatient@example.com`
   - Phone: `+998977777777`
   - Password: `patient123`

3. Нажать **Create User**
4. Увидеть сообщение "User created successfully!"

5. Перейти в **Patients** - увидеть нового пациента в списке

---

### Шаг 3: Создать медицинскую запись для пациента

1. Перейти в **Create Medical Data**
2. Заполнить форму:
   - Patient: выбрать `Тестовый Пациент (testpatient@example.com)`
   - Clinic ID: `1` (автоматически, т.к. врач из City Clinic)
   - Diagnosis: `Бронхит`
   - Complaint: `Кашель, температура 37.5`
   - Treatment: `Антибиотики, ингаляции`
   - Medications: `Амоксициллин 500мг 3 раза в день`
   - Allergies: `Нет`
   - Lab Results: `Общий анализ крови в норме`
   - Doctor Notes: `Повторный осмотр через неделю`

3. Нажать **Create Medical Record**
4. Увидеть "Medical record created successfully!"

5. Перейти в **Medical Records** - увидеть новую запись

---

### Шаг 4: Создать запрос на передачу данных

**Сценарий:** Пациент переезжает из City Clinic (1) в Regional Hospital (2)

1. Перейти в **Create Transfer**
2. Выбрать:
   - Medical Record: выбрать запись пациента `Бронхит`
   - From Clinic ID: `1` (автоматически)
   - To Clinic ID: `2` (Regional Hospital)
   - Patient User ID: автоматически подставляется

3. Нажать **Create Transfer Request**
4. Увидеть "Transfer request created! Patient must approve it."

5. Перейти в **Data Transfers** - увидеть запрос со статусом `pending`

---

### Шаг 5: Logout и войти как пациент

1. Нажать **Logout**
2. Войти:
   ```
   Email: testpatient@example.com
   Password: patient123
   ```

Вы видите:
- ✅ Dashboard пациента
- ✅ Только свои медицинские записи
- ✅ Pending transfer request

---

### Шаг 6: Пациент подтверждает передачу

1. Перейти в **Transfer Requests**
2. Увидеть запрос:
   - From Clinic: `Clinic #1`
   - To Clinic: `Clinic #2`
   - Status: `pending` (оранжевый)
   - Кнопки: **Approve** и **Reject**

3. Нажать **Approve**
4. Увидеть alert "Transfer approved successfully!"
5. Статус изменится на `confirmed` (зеленый)

---

### Шаг 7: Проверка в целевой клинике

1. Нажать **Logout**
2. Войти как врач Regional Hospital:
   ```
   Email: doctor2@hospital.com
   Password: doctor123
   ```

3. Перейти в **Medical Records**
4. ✅ **ВИДИМ**: Медицинская запись пациента появилась в клинике #2!
   - Patient ID совпадает
   - Diagnosis: `Бронхит`
   - Clinic ID теперь: `2`
   - Все данные скопировались

---

## 🔍 Что еще можно показать

### Admin возможности

Войти как admin:
```
Email: admin@system.com
Password: admin123
```

Видно:
- ✅ **All Patients** - все пациенты из всех клиник
- ✅ **All Employees** - все врачи со всех клиник
- ✅ **All Medical Data** - все медицинские записи
- ✅ **All Transfers** - все запросы на передачу данных
- ✅ Может создавать пользователей любых типов (patient, employee, admin)
- ✅ Может создавать данные для любых клиник

### Фильтрация по клиникам (Employee)

Войти как `doctor1@clinic.com`:
- Видит только пациентов и данные своей клиники (clinic_id = 1)

Войти как `doctor2@hospital.com`:
- Видит только пациентов и данные своей клиники (clinic_id = 2)

---

## 📊 Тестовые данные в системе

В базе уже есть:

**7 клиник:**
1. City Clinic
2. Regional Hospital
3. Medical Center Plus
4. Family Health Center
5. Children Hospital
6. Cardiology Center
7. Dental Clinic ProSmile

**18+ сотрудников** распределены по клиникам

**13+ пациентов** с медицинскими записями

**Множество медицинских записей** (ОРВИ, Гастрит, Аритмия, Кариес и т.д.)

**Несколько трансферов** в разных статусах (pending, confirmed, rejected)

---

## ✅ Checklist для демонстрации

- [ ] Бэкенд запущен (`docker compose ps` - все контейнеры UP)
- [ ] Фронтенд доступен (localhost:8000 открывается)
- [ ] Вход как врач работает
- [ ] Создание пациента работает
- [ ] Создание медданных работает
- [ ] Создание трансфера работает
- [ ] Вход пациентом работает
- [ ] Подтверждение трансфера работает
- [ ] Данные появляются в целевой клинике

---

## 🐛 Troubleshooting

**Проблема:** Docker не запускается
```bash
# Остановить все контейнеры
docker compose down -v

# Проверить порты
lsof -i :8099   # Должен быть свободен
lsof -i :50052  # Должен быть свободен
lsof -i :50053  # Должен быть свободен

# Запустить заново
docker compose up -d
```

**Проблема:** Форма не отправляется / ошибка API
- Открыть DevTools (F12)
- Проверить Console на ошибки
- Проверить Network tab - видны ли запросы к http://localhost:8099

**Проблема:** "Access-Control" ошибка
- Проверить что бэкенд запущен
- CORS уже настроен в core-service

---

## 🎓 Для защиты проекта

**Что говорить:**

1. "Это упрощенная система управления медицинскими данными для университетского проекта"

2. "Реализована микросервисная архитектура:
   - REST API Gateway (core-service)
   - gRPC сервис пользователей (user-service)
   - gRPC сервис медданных (data-transfer-service)"

3. "Три роли с разными правами:
   - Patient - видит только свои данные
   - Employee - видит данные своей клиники
   - Admin - видит всё"

4. "Полный цикл передачи данных:
   - Врач создает запрос
   - Пациент подтверждает
   - Данные копируются в целевую клинику"

5. "Технологии: Go, gRPC, PostgreSQL, JWT, bcrypt на бэке. Vanilla JS SPA на фронте"

**Что НЕ говорить:**
- ❌ "Это блокчейн система" (мы убрали блокчейн!)
- ❌ "Тут сложная криптография" (только JWT и bcrypt)
- ❌ "Загружаем файлы и изображения" (только строки!)
