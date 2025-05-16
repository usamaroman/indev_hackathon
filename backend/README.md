# CRM BACKEND SERVICE

[![Build Status][ci-badge]][ci-runs]
<!--- ![coverage][coverage-badge] --->

[ci-badge]:            https://github.com/tachki/crm-backend/actions/workflows/test.yaml/badge.svg
[ci-runs]:             https://github.com/tachki/crm-backend/actions
[coverage-badge]:      https://raw.githubusercontent.com/tachki/.github/badges/.badges/master/coverage.svg

Микросервис, с помощью которого бизнес сможет менеджить свой ассортимент и автоматизировать рутинные задачи

# Usage

Запустить сервис со всей инфраструктурой можно с помощью команды 

```shell
  make compose-up
# или
  docker compose up --build
```

Для запуска минимального бекенда  

```shell
  make dev
```

Документацию после запуска сервиса можно посмотреть по адресу `http://localhost:8080/swagger/index.html`
с портом 8080 по умолчанию


Перед запуском <u>интеграционных</u> тестов нужно явно указать путь до конфига и запустить docker и не выключать его до конца тестирования. При этом вы должны быть в <u>корне проекта</u>
```shell
  export CONFIG_PATH=$(pwd)/config.yaml && docker compose up --build
```

Запуск юнит тестов
```shell
  make unit-test
```

Запуск интеграционных тестов
```shell
  make integration-test
```

Запуск всех тестов
```shell
  make test
```

Запуск всех тестов с покрытием для получения отчёта в html
```shell
  make coverage-html
```
# Standart users
##### Администратор отеля (Admin)
```
  login: admin
  password: admin
  
  type: admin 
```
  
##### Пользователь (Customer)
```
  login: customer
  password: customer

  type: customer
```

# Configuration

Сконфигурировать приложение можно используя `config.yaml`, указав путь до файла в переменной `CONFIG_PATH`
#### example
```shell
  export CONFIG_PATH=config/config.yaml
```

```yaml
http:
  host: localhost
  port: 8080
postgresql:
  user: user
  password: password
  host: localhost
  port: 5432
  database: postgres
  ssl_mode: disable
  auto_create: false
jwt:
  sign_key: secret
  token_ttl: 60m
```

вместе с файлом приложение можно настроить используя перемнные окружения
- `CONFIG_PATH=path` - настройка расположения yaml конфиг файла; дефолт значение = "config.yaml"
- `APP_ENV=prod/dev` - настройка окружения приложения
- `LOG_LEVEL=debug/info/warn/error` - настройка уровня логирования; дефолт значение = "debug"

для всех полей из yaml файла есть переменные окружения для конфигурации

- `PORT`
- `HOST`
- `PG_USER`
- `PG_PASSWORD`
- `PG_HOST`
- `PG_PORT`
- `PG_DATABASE`
- `PG_SSL`
- `PG_AUTO_CREATE`
- `JWT_KEY`
- `JWT_TTL`

# Миграции БД 

```bash
docker compose down -v # только в самый первый раз для очистки данных, в следующий раз без флага -v
docker compose up postgres -d # запуск постгреса
docker compose up migrations --build # запускаем миграции со сборкой
make dev # запуск минимального бекенда
``` 
