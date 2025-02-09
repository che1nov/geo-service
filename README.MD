# Geo Service

Cервис для поиска и геокодирования адресов, реализованный с использованием принципов чистой архитектуры. Сервис обращается к [DaData API](https://dadata.ru/) для получения информации об адресах и поддерживает аутентификацию с использованием JWT. Основные функции сервиса:

- **Регистрация пользователя** – регистрация с хешированием пароля (bcrypt).
- **Авторизация пользователя (логин)** – выдача JWT-токена.
- **Поиск адресов** – эндпоинт `/api/address/search` принимает параметр `query` и возвращает список адресов.
- **Геокодирование** – эндпоинт `/api/address/geocode` принимает координаты `lat` и `lng` и возвращает информацию о местоположении.
- **Swagger-документация** – все эндпоинты задокументированы и доступны по адресу `/swagger/index.html`.

## Структура проекта

Проект организован согласно принципам чистой архитектуры и разделён на несколько слоёв:

- **internal/controller**  
  Содержит HTTP-обработчики (контроллеры), которые принимают запросы, выполняют базовую валидацию и вызывают методы бизнес-логики из слоя service.  
  Контроллеры используют интерфейс `Responder` для формирования ответов.

- **internal/service**  
  Содержит всю бизнес-логику (например, обращение к DaData API, обработку данных).

- **internal/models**  
  Содержит определения структур данных для запросов и ответов.

- **internal/repository**  
  Отвечает за взаимодействие с внешними сервисами (например, с DaData API).

- **internal/responder**  
  Определяет интерфейс и реализацию для формирования ответов (например, методы `RespondJSON` и `RespondError`).

## Требования

- **Go** 1.19 или выше (при локальной сборке).
- **Docker** и **Docker Compose** – для сборки и запуска приложения в контейнере.

## Быстрый старт с Docker Compose
Проект можно запустить с помощью Docker Compose.

## Makefile

В проекте присутствует Makefile с основными командами:
- make build: "Сборка бинарного файла"

- make run: "Запуск"

- make docker-build: "Сборка Docker-образа"

- make docker-up: "Запуск контейнеров через docker-compose"

- make clean: "Очистка сгенерированных файлов"


