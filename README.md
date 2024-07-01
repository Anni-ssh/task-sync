# TaskSync

TaskSync - это веб-приложение для управления временем выполнения задач.

## Функциональность

TaskSync предоставляет следующие основные функции:

### People

- **Создание пользователя**: Создание нового пользователя.
- **Получение списка пользователей**: Получение всех пользователей.
- **Получение информации о пользователе по ID**: Получение деталей пользователя по его ID.
- **Получение пользователей по фильтру**: Получение пользователей на основе заданных фильтров.
- **Обновление информации о пользователе**: Обновление данных существующего пользователя.
- **Удаление пользователя**: Удаление пользователя по его ID.

### Tasks

- **Создание задачи**: Создание новой задачи.
- **Получение задачи по ID**: Получение информации о задаче по её ID.
- **Получение списка задач**: Получение всех задач.
- **Обновление задачи**: Обновление данных существующей задачи.
- **Обновление пользователей в задаче**: Обновление пользователей, связанных с задачей.
- **Удаление задачи**: Удаление задачи по её ID.

### Time

- **Начало записи времени**: Начало записи времени для выполнения задачи.
- **Завершение записи времени**: Завершение записи времени для выполнения задачи.
- **Получение потраченного времени на задачи**: Получение времени, затраченного на выполнение задач определённым пользователем в заданном временном интервале.

## Использованные технологии

TaskSync разработан с использованием следующих технологий:

1. **Swagger**: Используется для документирования и взаимодействия с API.
2. **SQL Postgres**: Используется для хранения данных о задачах и пользователях.
3. **Chi Router**: Используется для маршрутизации HTTP запросов.
4. **Гексагональная архитектура**: Используется для организации кода и разделения бизнес-логики от инфраструктуры.
5. **Миграции**: Используются для управления изменениями в структуре базы данных.
6. **Логирование**: Используется slog для отслеживания действий и ошибок приложения.
7. **Конфигурация из .env файла**: Используется для загрузки конфигурационных параметров из файла .env.