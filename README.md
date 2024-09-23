# Private Keeper V2

**Private Keeper V2** — это выпускной проект Яндекс Практикума. Сервис состоит из серверной и клиентской части и предназначен для безопасного хранения данных в зашифрованном виде.

## Возможности сервиса

- Хранение данных в зашифрованном виде
- Шифрование данных ключом, который генерируется индивидуально для каждого пользователя
- Поддержка mTLS (Mutual TLS) между клиентом и сервером
- Пароли пользователей хранятся в виде хешей
- Кэширование ключей шифрования пользователя с использованием Redis
- Партицирование данных по видам сохраняемых данных
- Поддержка различных типов данных (например, текстовые данные, учетные данные, бинарные файлы и т.д.)

## Требования

Для запуска потребуется:

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Make](https://www.gnu.org/software/make/manual/make.html)

## Установка и запуск


1. **Сгенерируйте ключи сервера для mTLS:**

    ```bash
    make server-keys
    ```

2. **Сгенерируйте ключи клиента для mTLS:**

    ```bash
    make client-keys
    ```

### Запуск сервера

Выполните команду из корня проекта:

```bash
docker compose up
```

### Запуск клиента

Выполните команды из корня проекта:

```bash
make build
make run
```

## Базовое использование

1. **Запуск клиента:** Пользователь запускает клиентскую часть и может либо зарегистрироваться, либо войти в систему, если уже зарегистрирован.

2. **Сохранение данных:** Пользователь сохраняет необходимые данные в зашифрованном виде.

3. **Запрос данных:** При запросе данных можно использовать фильтры, чтобы избежать загрузки всех данных сразу.

4. **Вывод данных:** Пользователь может выбрать способ вывода данных — на экран терминала или сохранение в файл.

5. **Рабочая директория:** Пользователь может установить рабочую директорию, чтобы вводить только имя файла без полного пути.

6. **Сохранение бинарных данных:** Для сохранения бинарных данных необходимо установить рабочую директорию, куда будут сохраняться файлы.

## Технологии

- **mTLS**: защищенное соединение между клиентом и сервером
- **Шифрование**: индивидуальный ключ для каждого пользователя
- **Redis**: кэширование ключей шифрования
- **PostgreSQL**: база данных с партицированием данных

## Безопасность

- Данные на сервере хранятся в зашифрованном виде.
- Пароли пользователей защищены хешированием.
- Ключи шифрования пользователей кэшируются для оптимизации.

## Контакты

Если у вас есть вопросы или предложения, обратитесь к разработчику проекта через [GitHub](https://github.com/DenisKhanov/PrivateKeeperV2).

---

Private Keeper V2 — ваш надежный помощник для безопасного хранения данных!