-- Главные таблицы

-- Пользователи
CREATE TABLE IF NOT EXISTS  people_info (
    id SERIAL PRIMARY KEY,
    passport_series INT NOT NULL,
    passport_number INT NOT NULL,
    surname VARCHAR(50) NOT NULL,
    name VARCHAR(50) NOT NULL,
    patronymic VARCHAR(50),
    address TEXT NOT NULL,
    CONSTRAINT unique_passport UNIQUE (passport_series, passport_number)
);

-- Задания
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description TEXT
);

-- Контроль времени
CREATE TABLE IF NOT EXISTS time_entries (
    id SERIAL PRIMARY KEY,
    people_id INTEGER,
    task_id INTEGER NOT NULL UNIQUE,
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    duration INTERVAL GENERATED ALWAYS AS (end_time - start_time) STORED,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (people_id) REFERENCES people_info(id) ON DELETE SET NULL,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    CONSTRAINT chk_time CHECK (end_time >= start_time)
);

-- Индексы для ускорения поиска, замедляют добавление.
CREATE INDEX IF NOT EXISTS idx_people_info_passport ON people_info (passport_series, passport_number);
CREATE INDEX IF NOT EXISTS idx_time_entries_people_id ON time_entries (people_id);
CREATE INDEX IF NOT EXISTS idx_time_entries_task_id ON time_entries (task_id);


-- Функция для проверки длины паспортных данных 
CREATE OR REPLACE FUNCTION check_passport_length()
RETURNS TRIGGER AS $$
BEGIN
    -- Проверка длины passport_series
    IF NEW.passport_series IS NOT NULL AND (NEW.passport_series < 1000 OR NEW.passport_series > 9999) THEN
        RAISE EXCEPTION 'passport_series must be a 4-digit number';
    END IF;

    -- Проверка длины passport_number
    IF NEW.passport_number IS NOT NULL AND (NEW.passport_number < 100000 OR NEW.passport_number > 999999) THEN
        RAISE EXCEPTION 'passport_number must be a 6-digit number';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Создание триггера для проверки паспортных данных при вставке и обновлении
CREATE TRIGGER check_passport_length_trigger
BEFORE INSERT OR UPDATE ON people_info
FOR EACH ROW
EXECUTE FUNCTION check_passport_length();

-- Создание функции для вычисления duration для конкретной записи в time_entries по id
CREATE OR REPLACE FUNCTION calculate_duration_by_id(entry_id INTEGER)
RETURNS INTERVAL AS $$
DECLARE
    start_time TIMESTAMP;
    end_time TIMESTAMP;
    dur INTERVAL;
BEGIN
    -- Получаем start_time и end_time для данного entry_id
    SELECT te.start_time, te.end_time INTO start_time, end_time
    FROM time_entries te
    WHERE te.id = entry_id;

    -- Вычисляем duration
    dur := end_time - start_time;

    RETURN dur;
END;
$$ LANGUAGE plpgsql;

-- Создание хранимой процедуры для генерации рандомных значений.
-- Выполняется только 1 раз из-за уникальности значений task_id.
CREATE OR REPLACE PROCEDURE insert_test_data()
LANGUAGE plpgsql
AS $$
BEGIN
    -- Генерация рандомных данных для people_info
    INSERT INTO people_info (passport_series, passport_number, surname, name, patronymic, address)
    SELECT
        LPAD(FLOOR(RANDOM() * 10000)::TEXT, 4, '0'),  -- Генерация серии паспорта
        LPAD(FLOOR(RANDOM() * 1000000)::TEXT, 6, '0'),  -- Генерация номера паспорта
        CONCAT('Surname', FLOOR(RANDOM() * 1000)),  -- Генерация фамилии
        CONCAT('Name', FLOOR(RANDOM() * 1000)),  -- Генерация имени
        CONCAT('Patronymic', FLOOR(RANDOM() * 1000)),  -- Генерация отчества
        CONCAT('Address', FLOOR(RANDOM() * 1000))  -- Генерация адреса
    FROM generate_series(1, 50);  -- Вставка 50 случайных записей

    -- Генерация рандомных данных для tasks
    INSERT INTO tasks (title, description)
    SELECT
        CONCAT('Task ', FLOOR(RANDOM() * 1000)::TEXT),  -- Генерация названия задачи
        'Description for Task ' || FLOOR(RANDOM() * 1000)::TEXT  -- Генерация описания задачи
    FROM generate_series(1, 50); 

    -- Генерация данных для time_entries
    INSERT INTO time_entries (people_id, task_id, start_time, end_time)
    SELECT
        gs1 AS people_id,
        gs2 AS task_id,
        NOW(),  -- Дата начала работы
        NOW() + INTERVAL '1 day' * FLOOR(RANDOM() * 30)  -- Дата окончания работы
    FROM generate_series(1, 50) gs1
    JOIN generate_series(1, 50) gs2 ON gs1 = gs2  -- Соединение для уникальных пар (people_id, task_id) task_id - имеет constraint UNIQUE
    LIMIT 50;  

END;
$$;
