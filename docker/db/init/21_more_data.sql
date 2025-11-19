-- Дополнительные тестовые данные (больше пациентов и медданных)

DO $$
DECLARE
_hash TEXT := crypt('patient123', gen_salt('bf'));

    _p4_id BIGINT; _p5_id BIGINT; _p6_id BIGINT; _p7_id BIGINT; _p8_id BIGINT;
    _p9_id BIGINT; _p10_id BIGINT; _p11_id BIGINT; _p12_id BIGINT; _p13_id BIGINT;
    _p14_id BIGINT; _p15_id BIGINT;
BEGIN
    -------------------------------
    -- ПАЦИЕНТЫ 4–13 (как раньше)
    -------------------------------

    -- Пациент 4
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'patient4@example.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Александр', 'Лебедев', 'patient4@example.com', '+998988888881', 'patient', _hash, NULL)
        RETURNING id INTO _p4_id;
ELSE SELECT id INTO _p4_id FROM users WHERE email = 'patient4@example.com'; END IF;

    -- Пациент 5
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'patient5@example.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Ирина', 'Соловьева', 'patient5@example.com', '+998988888882', 'patient', _hash, NULL)
        RETURNING id INTO _p5_id;
ELSE SELECT id INTO _p5_id FROM users WHERE email = 'patient5@example.com'; END IF;

    -- Пациент 6
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'patient6@example.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Максим', 'Титов', 'patient6@example.com', '+998988888883', 'patient', _hash, NULL)
        RETURNING id INTO _p6_id;
ELSE SELECT id INTO _p6_id FROM users WHERE email = 'patient6@example.com'; END IF;

    -- Пациент 7
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'patient7@example.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Виктория', 'Григорьева', 'patient7@example.com', '+998988888884', 'patient', _hash, NULL)
        RETURNING id INTO _p7_id;
ELSE SELECT id INTO _p7_id FROM users WHERE email = 'patient7@example.com'; END IF;

    -- Пациент 8
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'patient8@example.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Денис', 'Романов', 'patient8@example.com', '+998988888885', 'patient', _hash, NULL)
        RETURNING id INTO _p8_id;
ELSE SELECT id INTO _p8_id FROM users WHERE email = 'patient8@example.com'; END IF;

    -- Пациент 9
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'patient9@example.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('София', 'Борисова', 'patient9@example.com', '+998988888886', 'patient', _hash, NULL)
        RETURNING id INTO _p9_id;
ELSE SELECT id INTO _p9_id FROM users WHERE email = 'patient9@example.com'; END IF;

    -- Пациент 10
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'patient10@example.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Артем', 'Захаров', 'patient10@example.com', '+998988888887', 'patient', _hash, NULL)
        RETURNING id INTO _p10_id;
ELSE SELECT id INTO _p10_id FROM users WHERE email = 'patient10@example.com'; END IF;

    -- Пациент 11
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'patient11@example.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Кристина', 'Макарова', 'patient11@example.com', '+998988888888', 'patient', _hash, NULL)
        RETURNING id INTO _p11_id;
ELSE SELECT id INTO _p11_id FROM users WHERE email = 'patient11@example.com'; END IF;

    -- Пациент 12
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'patient12@example.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Никита', 'Алексеев', 'patient12@example.com', '+998988888889', 'patient', _hash, NULL)
        RETURNING id INTO _p12_id;
ELSE SELECT id INTO _p12_id FROM users WHERE email = 'patient12@example.com'; END IF;

    -- Пациент 13
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'patient13@example.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Полина', 'Степанова', 'patient13@example.com', '+998988888890', 'patient', _hash, NULL)
        RETURNING id INTO _p13_id;
ELSE SELECT id INTO _p13_id FROM users WHERE email = 'patient13@example.com'; END IF;


    -------------------------------
    -- НОВЫЕ ПАЦИЕНТЫ 14 и 15
    -------------------------------

    -- Пациент 14
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'patient14@example.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Егор', 'Фадеев', 'patient14@example.com', '+998988888891', 'patient', _hash, NULL)
        RETURNING id INTO _p14_id;
ELSE SELECT id INTO _p14_id FROM users WHERE email = 'patient14@example.com'; END IF;

    -- Пациент 15
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'patient15@example.com') THEN
        INSERT INTO users (first_name, last_name, email, phone_number, type, password, clinic_id)
        VALUES ('Анна', 'Чернова', 'patient15@example.com', '+998988888892', 'patient', _hash, NULL)
        RETURNING id INTO _p15_id;
ELSE SELECT id INTO _p15_id FROM users WHERE email = 'patient15@example.com'; END IF;


    -----------------------------------------
    -- МЕДИЦИНСКИЕ ДАННЫЕ (включая новых)
    -----------------------------------------
    IF NOT EXISTS (SELECT 1 FROM medical_data WHERE user_id = _p4_id) THEN
        INSERT INTO medical_data (user_id, clinic_id, diagnosis, complaint, treatment, medications, allergies, doctor_notes, lab_results, created_at)
        VALUES
            -- Family Health Center (clinic 4)
            (_p4_id, 4, 'Бронхит', 'Кашель, затрудненное дыхание', 'Ингаляции, покой', 'Амброксол, Азитромицин', 'Нет', 'Контроль через неделю', 'Рентген: воспаление бронхов', NOW() - INTERVAL '4 days'),
            (_p5_id, 4, 'Мигрень', 'Сильная головная боль, светобоязнь', 'Покой, темнота', 'Суматриптан при приступе', 'Нет', 'Вести дневник приступов', 'МРТ: норма', NOW() - INTERVAL '12 days'),
            (_p6_id, 4, 'Остеохондроз', 'Боль в спине', 'ЛФК, массаж', 'НПВС курсом', 'Нет', 'Направление к неврологу', 'Рентген позвоночника: дегенеративные изменения', NOW() - INTERVAL '20 days'),

            -- Children Hospital (clinic 5)
            (_p7_id, 5, 'Ветряная оспа', 'Сыпь, температура', 'Карантин, симптоматическое лечение', 'Зеленка, жаропонижающее', 'Нет', 'Изоляция 10 дней', 'Клинический диагноз', NOW() - INTERVAL '6 days'),
            (_p8_id, 5, 'Кишечная инфекция', 'Диарея, рвота', 'Регидратация', 'Энтерофурил, Регидрон', 'Нет', 'Обильное питьё', 'Копрограмма: воспаление', NOW() - INTERVAL '2 days'),

            -- Cardiology Center (clinic 6)
            (_p9_id, 6, 'Ишемическая болезнь сердца', 'Боли за грудиной', 'Коррекция образа жизни', 'Аспирин, Аторвастатин', 'Нет', 'Наблюдение кардиолога', 'ЭКГ: ишемия', NOW() - INTERVAL '25 days'),
            (_p10_id, 6, 'Аритмия', 'Перебои в сердце', 'Медикаментозно', 'Конкор 5мг', 'Нет', 'Холтер через месяц', 'ЭКГ: экстрасистолия', NOW() - INTERVAL '8 days'),

            -- Dental Clinic (clinic 7)
            (_p11_id, 7, 'Кариес', 'Боль при жевании', 'Пломбирование', 'Анестезия местная', 'Лидокаин', 'Гигиена полости рта', 'Рентген: поражение дентина', NOW() - INTERVAL '1 day'),
            (_p12_id, 7, 'Пародонтит', 'Кровоточивость десен', 'Профчистка', 'Хлоргексидин', 'Нет', 'Контроль через 3 мес', 'Пародонтальные карманы', NOW() - INTERVAL '14 days'),

            -- Regional Hospital (clinic 2)
            (_p13_id, 2, 'Сахарный диабет 2 типа', 'Глюкоза 8.5', 'Диета, спорт', 'Метформин', 'Нет', 'Контроль сахара', 'HbA1c: 7.8%', NOW() - INTERVAL '30 days'),

            -- Новый пациент 14 (clinic 3)
            (_p14_id, 3, 'Гипертония', 'Головокружение, давление 150/95', 'Антигипертензивная терапия', 'Эналаприл 10 мг', 'Нет', 'Снижение соли, контроль АД', 'АД: 150/95, пульс 78', NOW() - INTERVAL '5 days'),

            -- Новый пациент 15 (clinic 5)
            (_p15_id, 5, 'Анемия', 'Усталость, слабость', 'Железосодержащие препараты', 'Феррум Лек', 'Нет', 'Пересдать анализы через месяц', 'Hb: 9.8, ферритин низкий', NOW() - INTERVAL '18 days');
END IF;


    -----------------------------------------
    -- ДОПОЛНИТЕЛЬНЫЕ TRANSFER
    -----------------------------------------

    IF NOT EXISTS (SELECT 1 FROM data_transfer WHERE user_id = _p4_id) THEN
        INSERT INTO data_transfer (medical_data_id, from_clinic_id, to_clinic_id, user_id, status, created_at)
SELECT md.id, 4, 2, _p4_id, 'pending', NOW() - INTERVAL '1 day'
FROM medical_data md WHERE md.user_id = _p4_id LIMIT 1;
END IF;

    IF NOT EXISTS (SELECT 1 FROM data_transfer WHERE user_id = _p9_id) THEN
        INSERT INTO data_transfer (medical_data_id, from_clinic_id, to_clinic_id, user_id, status, created_at)
SELECT md.id, 6, 3, _p9_id, 'confirmed', NOW() - INTERVAL '10 days'
FROM medical_data md WHERE md.user_id = _p9_id LIMIT 1;
END IF;

    IF NOT EXISTS (SELECT 1 FROM data_transfer WHERE user_id = _p7_id) THEN
        INSERT INTO data_transfer (medical_data_id, from_clinic_id, to_clinic_id, user_id, status, created_at)
SELECT md.id, 5, 4, _p7_id, 'rejected', NOW() - INTERVAL '15 days'
FROM medical_data md WHERE md.user_id = _p7_id LIMIT 1;
END IF;

    -- Новый трансфер 1
    IF NOT EXISTS (SELECT 1 FROM data_transfer WHERE user_id = _p14_id) THEN
        INSERT INTO data_transfer (medical_data_id, from_clinic_id, to_clinic_id, user_id, status, created_at)
SELECT md.id, 3, 6, _p14_id, 'pending', NOW() - INTERVAL '3 days'
FROM medical_data md WHERE md.user_id = _p14_id LIMIT 1;
END IF;

    -- Новый трансфер 2
    IF NOT EXISTS (SELECT 1 FROM data_transfer WHERE user_id = _p15_id) THEN
        INSERT INTO data_transfer (medical_data_id, from_clinic_id, to_clinic_id, user_id, status, created_at)
SELECT md.id, 5, 2, _p15_id, 'confirmed', NOW() - INTERVAL '7 days'
FROM medical_data md WHERE md.user_id = _p15_id LIMIT 1;
END IF;

END $$;
