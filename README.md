# WB Tech: level # 0 (Golang)	

## Тестовое задание

Необходимо разработать демонстрационный сервис с простейшим интерфейсом, отображающий данные о заказе. [Модель данных в формате JSON](/wildberries_l0/model.json) прилагается к заданию.

Что нужно сделать:

* Развернуть локально PostgreSQL
  + Создать свою БД
  + Настроить своего пользователя
  + Создать таблицы для хранения полученных данных
* Разработать сервис
  + Реализовать подключение и подписку на канал в nats-streaming
  + Полученные данные записывать в БД
  + Реализовать кэширование полученных данных в сервисе (сохранять in memory)
  + В случае падения сервиса необходимо восстанавливать кэш из БД
  + Запустить http-сервер и выдавать данные по id из кэша
* Разработать простейший интерфейс отображения полученных данных по id заказа

# Инструкция по запуску

1.
```shell
docker compose up -d
```
2.
```shell
cd pub-test
```
3.
```shell
go run .
``` 

```shell
docker compose down
```
# Работа:
___

Доступ к данным осуществляется по адресу
`localhost:8080/order/` по uid заказа


