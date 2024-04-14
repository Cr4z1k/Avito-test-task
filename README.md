# О проекте
Cервис, позволяющий показывать пользователям баннеры, в зависимости от требуемой фичи и тега пользователя, а также управлять баннерами и связанными с ними тегами и фичами.
### Запуск проекта
Склонируйте репозиторий:
```
https://github.com/Cr4z1k/Avito-test-task.git
```
### Конфигурация .env файла
Прежде чем запустить проект, убедитесь, что в вашей директории проекта создан файл `.env` с необходимыми переменными окружения. Далее представлен шаблон содержимого файла `.env`:
```
DB_PASS=
JWT_SALT=
DB_NAME=
DB_USER=
```
### Запуск сервера
1. Если у вас UNIX-подобная ОС, такая как Linux или macOS, выполните следующие команды:
   - Сборка проекта: `make build`
   - После сборки запустите проект:
     - Обычный запуск: `make run`
     - Запуск с тестами: `make run-w-tests`
     - Обычный запуск на заднем фоне: `make run-bg`
     - Запуск с тестами на заднем фоне: `run-w-tests-bg`
   - Остановка/Удаление контейнеров:
     - Остановка: `make stop`
     - Удаление: `make down`

2. Если в вашей OC не установлена утилита make, то можно выполнить следующие команды:
   - Сборка проекта: `docker-compose build`
   - После сборки запустите проект:
     - Обычный запуск: `docker-compose up`
     - Обычный запуск на заднем фоне: `docker-compose up -d`
     - Запуск с тестами:
       - bash:
         - Обычный: `RUN_TESTS=true docker-compose up`
         - На заднем фоне: `RUN_TESTS=true docker-compose up -d`
       - powershell:
         - Обычный: `$env:RUN_TESTS = "true"; docker-compose up`
         - На заднем фоне: `$env:RUN_TESTS = "true"; docker-compose up -d`
   - Остановка/Удаление контейнеров:
     - Остановка: `docker-compose stop`
     - Удаление: `docker-compose down`
# Информация для проверяющего
Все запросы выполняются по базовому пути `localhost:8080` при запуске проекта на своем компьютере, при запуске проекта с помощью docker-compose - `localhost:8000`.
В БД для облегечения работы заранее создаются 3 фичи с id(1, 2, 3) и 5 тегов с id(1, 2, 3, 4, 5). Причины такого подхода, а также работа endpoint'ов для добавления фичей и тегов описана в разделе `"Дополнительная информация"`, там же описан endpoint для получения токенов.
### Endpoints
Далее указаны примеры маршрутов при запуске через docker-compose:
- User
  - GET `http://localhost:8000/user_banner?tag_id=2&feature_id=1`
- Admin
  - GET `http://localhost:8000/user_banner?tag_id=2&feature_id=1`
  - GET `http://localhost:8000/banner?tag_id=1&feature_id=1&limit=100&offset=0`
  - POST `http://localhost:8000/banner`
    - JSON:
      ```json
      {
        "tag_ids": [
          1,
          2
        ],
        "feature_id": 1,
        "content": {
          "title": "some_title",
          "text": "some_text",
          "url": "some_url"
        },
        "is_active": true
      }
      ```
  - PATCH `http://localhost:8000/banner/1`
    - JSON:
      ```json
      {
        "tag_ids": [
          1,
          2
        ],
        "feature_id": 1,
        "content": {
          "title": "some_title",
          "text": "some_text",
          "url": "some_url"
        },
        "is_active": true
      }
      ```
  - DELETE `http://localhost:8000/banner/1`
### Дополнительная информация
В сервис также было добавлено несколько endpoint'ов для получения токенов и добавления в БД фичей и тегов. Автор предполагает, что функционал реализации токенов, фичей и тегов возлагается на другие сервисы, поэтому в данном сервисе для них был воссоздан минимальный функционал, дабы продемонстрировать его работоспособность.
1. Для получения JWT токенов используется endpoint "/get_token/:value" с методом GET. Для получения токена администратора в значении параметра value требуется указать 1. Любое другое значение генерирует токен для обычного пользователя.
2. Для добавления фичей в БД используется endpoint "/feature" с методом POST. В теле запроса необходимо передать JSON файл типа :
    ```json
    {
        "features": [
            1,
            2
        ]
    }
    ```
3. Для добавления тегов в БД используется endpoint "/tag" с методом POST. В теле запроса необходимо передать JSON файл типа :
    ```json
    {
        "tags": [
            1,
            2
        ]
    }
    ```
Т.к. был воссоздан минимальный функционал, то endpoint'ы для добавления фичей и тегов будут иметь статус ошибки сервера(500) при попытке добавить уже существующий фичу/тег вместо статуса конфлитка(409). При успешном добавлении статус-код запроса будет равен 201, а при неправильном формате JSON файла или пустом слайсе - 400.