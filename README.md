[![Go Report Card](https://goreportcard.com/badge/github.com/anfilat/final-stats)](https://goreportcard.com/report/github.com/anfilat/final-stats)

# Финальный проект по курсу [Golang Developer от Отус](https://otus.ru/lessons/golang-professional/) - Приложение осуществляющее системный мониторинг

Задание на проект находится в папке docs

## Что делает?

Собирает метрики о системе и отправляет их клиентам по gRPC.
Клиент при запросе передает параметры N и M. Приложение отправляет метрики каждые N секунд,
усредняя их за последние M секунд. Работает на Linux (Ubuntu) и Windows.

### Метрики

- Средняя загрузка системы
- Средняя загрузка CPU
- Загрузка дисков
- Информация о дисках по каждой файловой системе

## Внутреннее устройство

Приложение состоит из:

- gRPC сервер. Принимает запросы и в поточном режиме отдает метрики подключившимся клиентам
- Сервис клиентов. Хранит список подключенных клиентов и в соответствии с параметрами клиента отсылает ему метрики каждые N секунд
- Сервис сбора метрик. Каждую секунду запрашивает метрики у коллекторов, ответственных за их получение,
и отправляет все накопленные посекундные метрики сервису клиентов.
- Коллекторы, ответственные за сбор конкретных метрик. Каждый выполняется в своем потоке

Сервис сбора метрик работает постоянно, собирая метрики в памяти. Устаревшие данные (время устаревания задается в конфиге) удаляются.
Сервисы клиентов и сбора метрик общаются через канал. По нему сервису клиентов каждую секунду передается полная копия собранных метрик.
Сервис клиентов передает каждому клиенту при подключении канал, по которому будут приходить метрики и функцию отключения.
