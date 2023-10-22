# mincer / мясорубка

Клиент-серверная игра.

## Интерфейс

![img.png](img.png)

### HUD

- Красная шкала со здоровьем
- Cool down шкала скорости атак
- Небольшое синее окно с классом персонажа, оружием и уроном
- Справа информационные сообщения (убийства, присоединения/отсоединения других игроков)
- Снизу координаты и FPS

### Динамика игры

- Зелёное поле - поле битвы
- Игровой персонаж обведён белой обводкой.
- Синие персонажи - маги, красные - воины, зелёные - лучники, серые - мёртвые игроки.
- Маги и лучники стреляют фаерболами (маленькие красные кружки) и стрелами (черные точки), соответственно.
- Воины подходят и наносят урон вблизи.
- У воинов и магов оружие наносит урон по радиусу. Лучники - только в одну цель.
- В игру добавлено несколько ботов.

## Сборка и запуск

### Подготовка

1. Установить кодогенератор для proto файлов:
    ```shell
    sudo apt install protobuf-compiler protoc-gen-go
    ```

2. Установить утилиту make для сборки:
    ```shell
    sudo apt-get install make
    ```

3. Установить C компилятор для [Mac](https://ebitengine.org/en/documents/install.html?os=darwin#Installing_a_C_compiler) или для [Linux](https://ebitengine.org/en/documents/install.html?os=linux#Installing_a_C_compiler) и зависимости для [Linux](https://ebitengine.org/en/documents/install.html?os=linux#Installing_dependencies), чтобы собрать клиентское приложение.

### Сборка серверной и клиентской части

```shell
make build
```

Бинарные файлы будут помещены в [/bin](/bin).

### Генерация конфигурационных файлов

1. Сгенерировать пароль для Postgres:
   ```shell
   printf $(head -c 32 /dev/random | base64) > ./server/config/pg_password
   ```

2. Сгенерировать приватный ключ для взаимодействия netcode с клиентами:
   ```shell
   openssl rand -base64 32 > ./server/config/nc_private_key
   ```

3. Сгенерировать конфиг для миграций:
   ```shell
   printf '[database]
   host = postgres
   port = 5432
   database = mincer
   user = mincer
   password = '$(cat ./server/config/pg_password)'
   
   [data]
   ' > ./server/config/migrations.conf
   ```

### Запуск

1. Запустить базу данных в докере:
   ```shell
   docker-compose up --build -d postgres
   ```

2. Провести миграции в БД:
   ```shell
   docker-compose up --build migrations
   ```

3. Запустить сервер:
   ```shell
   POSTGRES_PASSWORD_FILE=server/config/pg_password \
   NC_PRIVATE_KEY_FILE=server/config/nc_private_key \
   LOG_LEVEL=info ./bin/server
   ```

4. Запустить клиент (или несколько клиентов):
   ```shell
   ./bin/client -a http://localhost:8080/token
   ```
   localhost или локальный IP сервера

## Управление

### На сервере

- Ctrl+C для изящного выключения
- переменные среды (см. [env.go](server/internal/configuration/env.go))

### На клиенте

- WASD или стрелки для движения персонажа
- Space для удара или запуска фаербола/стрелы
- R для возрождения
- Esc для выхода из игры
