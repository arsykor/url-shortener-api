# REST API сервис для сокращения ссылок

Основной функционал сервиса:

* сохранение короткого представления заданного URL
* redirect на исходный URL при переходе по сохраненному ранее короткому представлению

 
При сохранении URL происходит автоматическая валидация.
В качестве хранилища используется PostgreSQL или in-memory.

## Начало работы
1. Клонирование репозитория GitHub:

```shell
# Clone from GitHub
git clone https://github.com/arsykor/url-shortener-api
```
Параметры подключения к базе данных и сервера можно изменить в файле `config.yml`.

2. Настройка подключения к PostgreSQL



Данный пункт можно пропустить в случае использования памяти приложения в качестве хранилища.

Запросы для создания БД, необходимых таблиц и функций находятся в файле  `database.sql`.

2. Запуск приложения:

```shell
# Run
go run ./...
```
По умолчанию в проекте будет использоваться память приложения в качестве хранилища. Для подключения PostgreSQL необходимо присвоить значение "db" string параметру "storage" при запуске приложения.

```shell
# Help
go run ./... --help

# PostgreSQL as a storage
go run ./... --storage db 
```

## API эндпоинты



* `POST /` - сохранение оригинального URL в хранилище и возвращение сокращённого

* `GET /task/:id` - redirect на соответствующий исходный URL при передаче сгенерированной ссылки 

Для тестирования корректной работы приложения можно отправить следующий запрос в терминале:

```shell
# store URL: POST /task/create

curl -X POST -H "Content-Type: application/json" -d '
[
  {
    "url": "https://job.ozon.ru/fintech/"
  }
]
' http://localhost:8080
```
Ответ должен выглядеть примерно так:
```json
{
    "link": "http://localhost:8080/vgzJmQi5K5"
}
```
Ошибки возвращаются в стандартном формате:

```json
{
    "code": 400,
    "message": "the link for URL https://job.ozon.ru/fintech/ already exists, try using ID = vgzJmQi5K5"
}
```
