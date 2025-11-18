-- Дополнительные тестовые данные (больше пациентов и медданных)

DO $$
DECLARE
    _hash TEXT := crypt('patient123', gen_salt('bf'));
    _p4_id BIGINT; _p5_id BIGINT; _p6_id BIGINT; _p7_id BIGINT; _p8_id BIGINT;
    _p9_id BIGINT; _p10_id BIGINT; _p11_id BIGINT; _p12_id BIGINT; _p13_id BIGINT;
BEGIN
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

    -- Медданные для новых пациентов
    IF NOT EXISTS (SELECT 1 FROM medical_data WHERE user_id = _p4_id) THEN
        INSERT INTO medical_data (user_id, clinic_id, diagnosis, complaint, treatment, medications, allergies, doctor_notes, lab_results, created_at)
        VALUES
            -- Family Health Center (clinic 4)
            (_p4_id, 4, 'Бронхит', 'Кашель, затрудненное дыхание', 'Ингаляции, покой', 'Амброксол, Азитромицин', 'Нет', 'Контроль через неделю', 'Рентген: воспаление бронхов', NOW() - INTERVAL '4 days'),
            (_p5_id, 4, 'Мигрень', 'Сильная головная боль, светобоязнь', 'Покой, темнота', 'Суматриптан при приступе', 'Нет', 'Вести дневник приступов', 'МРТ: норма', NOW() - INTERVAL '12 days'),
            (_p6_id, 4, 'Остеохондроз', 'Боль в спине', 'ЛФК, массаж', 'НПВС курсом', 'Нет', 'Направление к неврологу', 'Рентген позвоночника: дегенеративные изменения', NOW() - INTERVAL '20 days'),

            -- Children Hospital (clinic 5)
            (_p7_id, 5, 'Ветряная оспа', 'Сыпь, температура', 'Карантин, симптоматическое лечение', 'Зеленка местно, жаропонижающее', 'Нет', 'Изоляция 10 дней', 'Клинический диагноз', NOW() - INTERVAL '6 days'),
            (_p8_id, 5, 'Кишечная инфекция', 'Диарея, рвота', 'Регидратация', 'Энтерофурил, Регидрон', 'Нет', 'Обильное питье', 'Копрограмма: воспаление', NOW() - INTERVAL '2 days'),

            -- Cardiology Center (clinic 6)
            (_p9_id, 6, 'Ишемическая болезнь сердца', 'Боли за грудиной', 'Коррекция образа жизни', 'Аспирин, Аторвастатин, Бета-блокаторы', 'Нет', 'Наблюдение кардиолога', 'ЭКГ: признаки ишемии, ЭхоКГ: снижение фракции выброса', NOW() - INTERVAL '25 days'),
            (_p10_id, 6, 'Аритмия', 'Перебои в работе сердца', 'Медикаментозная терапия', 'Конкор 5мг утром', 'Нет', 'Холтер мониторинг через месяц', 'ЭКГ: экстрасистолия', NOW() - INTERVAL '8 days'),

            -- Dental Clinic (clinic 7)
            (_p11_id, 7, 'Кариес', 'Боль при жевании', 'Пломбирование', 'Анестезия местная', 'Лидокаин', 'Гигиена полости рта', 'Рентген: поражение дентина', NOW() - INTERVAL '1 day'),
            (_p12_id, 7, 'Пародонтит', 'Кровоточивость десен', 'Профчистка, лечение', 'Хлоргексидин полоскание', 'Нет', 'Контроль через 3 месяца', 'Пародонтальные карманы 4-5мм', NOW() - INTERVAL '14 days'),

            -- Regional Hospital (clinic 2) - дополнительные
            (_p13_id, 2, 'Сахарный диабет 2 типа', 'Повышенный сахар 8.5', 'Диета, физнагрузки', 'Метформин 1000мг 2 раза в день', 'Нет', 'Контроль гликемии ежедневно', 'HbA1c: 7.8%, глюкоза натощак: 8.5', NOW() - INTERVAL '30 days');
    END IF;

    -- Дополнительные трансферы
    IF NOT EXISTS (SELECT 1 FROM data_transfer WHERE user_id = _p4_id) THEN
        INSERT INTO data_transfer (medical_data_id, from_clinic_id, to_clinic_id, user_id, status, created_at)
        SELECT md.id, 4, 2, _p4_id, 'pending', NOW() - INTERVAL '1 day'
        FROM medical_data md WHERE md.user_id = _p4_id AND md.clinic_id = 4 LIMIT 1;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM data_transfer WHERE user_id = _p9_id) THEN
        INSERT INTO data_transfer (medical_data_id, from_clinic_id, to_clinic_id, user_id, status, created_at)
        SELECT md.id, 6, 3, _p9_id, 'confirmed', NOW() - INTERVAL '10 days'
        FROM medical_data md WHERE md.user_id = _p9_id AND md.clinic_id = 6 LIMIT 1;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM data_transfer WHERE user_id = _p7_id) THEN
        INSERT INTO data_transfer (medical_data_id, from_clinic_id, to_clinic_id, user_id, status, created_at)
        SELECT md.id, 5, 4, _p7_id, 'rejected', NOW() - INTERVAL '15 days'
        FROM medical_data md WHERE md.user_id = _p7_id AND md.clinic_id = 5 LIMIT 1;
    END IF;

END $$;
