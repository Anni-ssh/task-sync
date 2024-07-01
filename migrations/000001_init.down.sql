-- Удаление таблиц
DROP TABLE IF EXISTS time_entries;
DROP TABLE IF EXISTS people_info;
DROP TABLE IF EXISTS tasks;

-- Удаление хранимых процедур
DROP PROCEDURE IF EXISTS insert_test_data;

-- Удаление триггеров
DROP TRIGGER IF EXISTS check_length_trigger ON people_info;

-- Удаление функций
DROP FUNCTION IF EXISTS calculate_duration_by_id;
DROP FUNCTION IF EXISTS check_length() CASCADE;

-- Удаление индексов
DROP INDEX IF EXISTS idx_people_info_passport;
DROP INDEX IF EXISTS idx_time_entries_people_id;
DROP INDEX IF EXISTS idx_time_entries_task_id;
