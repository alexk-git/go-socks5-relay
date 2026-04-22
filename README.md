# SOCKS5 Proxy Server on Go

[![Go Version](https://img.shields.io/badge/Go-1.22+-blue.svg)](https://golang.org)
[![Systemd](https://img.shields.io/badge/systemd-ready-green.svg)](https://systemd.io)

Производственно-готовый SOCKS5 прокси-сервер с аутентификацией, написанный на Go. Поддерживает graceful shutdown, гибкую настройку через флаги и переменные окружения, а также установку как systemd сервис.

Быстрый старт:
[Установка Go](#-установка-go)
[Сборка и запуск](#-сборка-и-запуск)
[Docker Compose](#-docker-compose)

---

##  Оглавление

- [Возможности](#-возможности)
- [Требования](#-требования)
- [Настройка конфигурации](#-настройка-конфигурации)
- [Сборка и запуск](#-сборка-и-запуск)
- [Использование](#-использование)
- [Установка как systemd сервиса](#-установка-как-systemd-сервиса)
- [Docker](#-docker)
- [Структура проекта](#-структура-проекта)
- [Разработка](#-разработка)
- [Устранение проблем](#-устранение-проблем)
- [Для LLM (подсказка)](#-для-llm-подсказка)

---

##  Возможности

-  **SOCKS5 протокол** с поддержкой TCP и UDP (RFC 1928)
-  **Аутентификация** Username/Password (RFC 1929)
-  **Graceful shutdown**  корректное завершение всех соединений
-  **Гибкая настройка**  файл конфигурации, флаги командной строки, переменные окружения
-  **Уровни логирования**  error, warn, info, debug
-  **Фильтрация шумных ошибок**  без спама в логах (EOF, broken pipe и т.д.)
-  **Логирование подключений**  видно все новые соединения и их закрытие
-  **Systemd интеграция**  можно установить как сервис
-  **Docker support**  готовый Dockerfile
-  **Чистая архитектура**  разделение на пакеты

---

##  Требования

- **Go 1.22+** (для сборки из исходников)
- **Make** (опционально, для удобства)
- **Git** (для скачивания зависимостей)
- **Docker** (опционально, для контейнеризации)

### Установка Go (Linux, официальный архив)

Актуальную версию и ссылку для скачивания смотрите на [go.dev/dl](https://go.dev/dl/).

```bash
# Скачать архив
wget https://go.dev/dl/go1.22.5.linux-amd64.tar.gz

# Распаковать в /usr/local
sudo tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz

# Добавить Go в PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Проверить установку
go version
```

> Для ARM-систем (например, Raspberry Pi) используйте архив `linux-arm64.tar.gz` вместо `linux-amd64.tar.gz`.

---

##  Настройка конфигурации

### Файл конфигурации (`config/.env`)

Пример конфигурационного файла `config/.env.example`:

```
# IP адрес для прослушивания
# 0.0.0.0 - слушать на всех интерфейсах (доступно извне)
# 127.0.0.1 - только локальный доступ
ip=0.0.0.0

# Порт для прослушивания
port=1080

# Имя пользователя для аутентификации
username=proxyuser

# Пароль для аутентификации
password=CHANGE_THIS_STRONG_PASSWORD
```

На его основе можно сделать свой собственный конфиг-файл `config/.env` или сгенерировать с помощью скрипта `scripts/generate-env.sh`.

### Приоритет настроек

Настройки применяются в следующем порядке (от высшего приоритета к низшему):

| Приоритет | Источник | Пример |
|-----------|----------|--------|
| 1 | Флаги командной строки | `-port 8080` |
| 2 | Переменные окружения | `SOCKS5_CONFIG=/path/to/config` |
| 3 | Файл конфигурации | `.env` |
| 4 | Значения по умолчанию | `port=1080` |

### Описание параметров

| Параметр | Описание | Значение по умолчанию |
|----------|----------|---------------------|
| `ip` | IP адрес для прослушивания. `0.0.0.0` означает все интерфейсы | `0.0.0.0` |
| `port` | Порт для прослушивания | `1080` |
| `username` | Имя пользователя для аутентификации | (обязательный параметр) |
| `password` | Пароль для аутентификации | (обязательный параметр) |

### Переопределение через переменные окружения

```
# Указать путь к конфигурационному файлу
export SOCKS5_CONFIG=/etc/my-proxy/.env

# Включить debug режим
export SOCKS5_DEBUG=1

# Установить уровень логирования
export SOCKS5_LOG_LEVEL=debug

# Запуск с переменными окружения
./build/socks5-proxy
```

### Переопределение через флаги командной строки

```
# Переопределить порт
./build/socks5-proxy -port 8080

# Переопределить IP
./build/socks5-proxy -ip 127.0.0.1

# Использовать кастомный конфиг
./build/socks5-proxy -config ./custom.properties
```

---

##  Сборка и запуск

### Предварительные требования

Убедитесь, что установлены необходимые инструменты:

```
# Проверка версии Go (требуется 1.22+)
go version

# Проверка Make
make --version

# Проверка Git
git --version
```

### Установка зависимостей
### Автоматическая генерация конфигурации

При первом запуске прокси автоматически создаст файл `.env` с случайными учетными данными.
Для этого в Makefile добавлены специальные команды:

```bash
# Сгенерировать .env файл с случайным именем пользователя и паролем
make generate-env

# Или выполнить полную настройку
make setup
```

При использовании `make run`, `make run-debug` или `make run-dev` генерация `.env` выполняется автоматически, если файл отсутствует.

Сгенерированные учетные данные выводятся в консоль. Сохраните их для подключения к прокси.

Пример вывода:
```
Generated .env with:
  Username: user_8993ceb4
  Password: b19c96874703a456091fbba2
Please save these credentials for connecting to the proxy.
```

Если вы хотите использовать свои учетные данные, просто отредактируйте файл `.env` вручную.

### Сборка проекта

#### Через Makefile (рекомендуется)

```
# Стандартная сборка
make build

# Сборка с очисткой перед сборкой
make clean && make build
```

#### Вручную

```
# Сборка бинарника
go build -o build/socks5-proxy ./cmd/socks5-proxy
```

После успешной сборки бинарник появится в директории `build/`:

```
ls -la build/
# socks5-proxy
```

### Запуск

#### Через Makefile

```
# Собрать и запустить
make run

# Собрать и запустить с debug режимом
make run-debug

# Запуск в development режиме (полное логирование)
make run-dev
```

#### Напрямую

```
# Обычный запуск
./build/socks5-proxy

# С debug режимом
./build/socks5-proxy -debug

# С указанием уровня логирования
./build/socks5-proxy -log-level debug

# С переопределением порта
./build/socks5-proxy -port 8080
```

### Проверка работы

После запуска вы должны увидеть:

```
[socks5] [I]: === SOCKS5 Proxy Server ===
[socks5] [I]: Адрес: 0.0.0.0:1080
[socks5] [I]: Аутентификация: Username/Password
[socks5] [I]: Пользователь: "proxyuser"
[socks5] [I]: Конфигурация: .env
[socks5] [I]: Уровень логирования: info
[socks5] [I]: Debug режим: false
[socks5] [I]: ===========================
[socks5] [I]: SOCKS5-прокси запущен на 0.0.0.0:1080 (TCP + UDP)
```

В другом терминале проверьте работу прокси:

```
# Замените username:password на ваши данные из .env
curl --proxy socks5://proxyuser:proxypass@localhost:1080 https://api.ipify.org

# Должен вернуться ваш IP-адрес (хоста, на котором запущен proxy)
```

### Остановка сервера

Нажмите `Ctrl+C` для graceful shutdown. Сервер корректно завершит все активные соединения перед остановкой.

---

##  Использование

### Флаги командной строки

```
./build/socks5-proxy [options]
```

| Флаг | Сокращение | Описание | Пример |
|------|------------|----------|--------|
| `-config` | - | Путь к файлу конфигурации | `-config ./custom.properties` |
| `-port` | - | Переопределить порт | `-port 8080` |
| `-ip` | - | Переопределить IP | `-ip 127.0.0.1` |
| `-log-level` | - | Уровень логирования (error/warn/info/debug) | `-log-level debug` |
| `-debug` | `-d` | Включить debug режим | `-debug` или `-d` |
| `-help` | `-h` | Показать справку | `-help` |

### Примеры запуска

```
# Обычный запуск
./build/socks5-proxy

# С debug режимом (видно все подключения)
./build/socks5-proxy -debug

# С указанием уровня логирования
./build/socks5-proxy -log-level debug

# Переопределение порта
./build/socks5-proxy -port 8080

# Только локальный доступ
./build/socks5-proxy -ip 127.0.0.1

# Кастомный конфиг
./build/socks5-proxy -config /etc/my-proxy/.env

# Комбинация
./build/socks5-proxy -debug -port 1080 -log-level debug
```

### Переменные окружения

```
# Включить debug режим
export SOCKS5_DEBUG=1

# Указать путь к конфигурации
export SOCKS5_CONFIG=/etc/proxy/.env

# Установить уровень логирования
export SOCKS5_LOG_LEVEL=debug

# Запустить
./build/socks5-proxy
```

---

##  Установка как systemd сервиса

### Автоматическая установка

```
# Собрать и установить
make install

# Настроить конфигурацию
sudo nano /etc/socks5-proxy/.env

# Запустить сервис
sudo systemctl enable socks5-proxy
sudo systemctl start socks5-proxy

# Проверить статус
sudo systemctl status socks5-proxy
```

### Ручная установка

```
# Создать директории
sudo mkdir -p /etc/socks5-proxy /var/log/socks5-proxy

# Скопировать бинарник
sudo cp build/socks5-proxy /usr/local/bin/
sudo chmod +x /usr/local/bin/socks5-proxy

# Скопировать конфигурацию
sudo cp .env /etc/socks5-proxy/
sudo chmod 600 /etc/socks5-proxy/.env

# Создать пользователя (опционально)
sudo useradd -r -s /bin/false -M -d /etc/socks5-proxy socks5
sudo chown -R socks5:socks5 /etc/socks5-proxy /var/log/socks5-proxy

# Скопировать файл сервиса
sudo cp socks5-proxy.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable socks5-proxy
sudo systemctl start socks5-proxy
```

### Файл сервиса (`socks5-proxy.service`)

```
[Unit]
Description=SOCKS5 Proxy Server
After=network.target

[Service]
Type=simple
User=nobody
Group=nogroup
ExecStart=/usr/local/bin/socks5-proxy -config /etc/socks5-proxy/.env
Restart=on-failure
RestartSec=10
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
```

### Управление сервисом

```
# Статус
sudo systemctl status socks5-proxy

# Запуск/остановка
sudo systemctl start socks5-proxy
sudo systemctl stop socks5-proxy
sudo systemctl restart socks5-proxy

# Просмотр логов
sudo journalctl -u socks5-proxy -f
sudo journalctl -u socks5-proxy -n 50
```

---

##  Docker

### Сборка образа

```
docker build -t socks5-proxy .
```

### Запуск контейнера

```
# Обычный запуск
docker run -p 1080:1080 socks5-proxy

# С debug режимом
docker run -p 1080:1080 socks5-proxy -debug

# С кастомным конфигом
docker run -p 1080:1080 -v $(pwd)/.env:/app/.env socks5-proxy

# В фоновом режиме
docker run -d -p 1080:1080 --name socks5-proxy socks5-proxy -debug
```

При первом запуске контейнера автоматически генерируется файл `.env` со случайными учетными данными. Для просмотра сгенерированных учетных данных проверьте логи контейнера:

```bash
docker logs <container_name>
```

Если необходимо использовать свои учетные данные, смонтируйте свой файл `.env` в контейнер:

```bash
docker run -p 1080:1080 -v $(pwd)/.env:/app/.env socks5-proxy
```

### Docker Compose

В проекте присутствует файл `docker-compose.yml` для удобного запуска через Docker Compose.

#### Запуск

```bash
# Запуск в фоновом режиме
docker-compose up -d

# Просмотр логов
docker-compose logs -f

# Остановка
docker-compose down
```

#### Работа с `.env` файлом

Конфигурация SOCKS5 прокси хранится в файле `.env`. Проект использует следующую логику:

1. **Автоматическая генерация**: При запуске контейнера через Docker Compose, если файл `config/.env` отсутствует, он будет автоматически сгенерирован на основе шаблона `config/.env.example` (или корневого `.env.example`, если шаблон в config отсутствует).

2. **Использование существующего файла**: Если `config/.env` уже существует, он будет использован без изменений. Это позволяет вам настроить параметры прокси (IP, порт, учётные данные) и сохранить их между перезапусками.

3. **Расположение файла**: По умолчанию `.env` файл находится в директории `config/` на хосте, которая монтируется в контейнер как `/app/config`. Это обеспечивает сохранение конфигурации между перезапусками контейнера.

Для просмотра сгенерированных параметров проверьте логи контейнера:
```bash
docker-compose logs socks5-proxy | grep -A2 "Configuration loaded"
```

**Важно**: Если вы измените шаблон `.env.example`, существующий `.env` файл не обновится автоматически. Для принудительной генерации удалите `config/.env` и перезапустите контейнер.

#### Сохранение `.env` между перезапусками


По умолчанию `.env` файл сохраняется в директории `config/` на хосте благодаря настроенному тому `./config:/app/config` в `docker-compose.yml`.

#### Использование своего `.env` файла

Если вы хотите использовать свой файл конфигурации, смонтируйте его в контейнер:

```yaml
volumes:
  # Для использования своего .env файла
  - ./.env:/app/.env:ro
```
 
Или передайте путь через переменную окружения:

```yaml
environment:
  SOCKS5_CONFIG: /app/.env
```

#### Настройка портов и других параметров

Порт прокси настраивается в файле `config/.env` (ключ `port`). Если файл `config/.env` отсутствует, он будет автоматически сгенерирован из `config/.env.example` при запуске контейнера.

Для автоматического обновления маппинга портов в Docker Compose используйте скрипт `generate-compose-env.sh`, который создаёт файл `.env` в корне проекта с переменной `SOCKS5_PORT`. Этот файл используется docker-compose для подстановки порта в маппинге.

Вы можете запустить генерацию и запуск контейнера одной командой:

```bash
make docker-up
```

Вручную:

```bash
./scripts/generate-compose-env.sh
docker-compose up
```

Если вы хотите изменить порт, отредактируйте `config/.env` (или `config/.env.example`, если `config/.env` не существует) и перезапустите контейнер.

Для включения debug режима раскомментируйте переменные окружения:

```yaml
environment:
  SOCKS5_DEBUG: "true"
  SOCKS5_LOG_LEVEL: "debug"
```
##  Структура проекта

```
go-socks5-relay/
 cmd/
    socks5-proxy/
        main.go                 # Точка входа, обработка флагов
 internal/
    config/
       config.go               # Структура Config и методы
       loader.go               # Загрузка конфигурации из файла
    logger/
       logger.go               # Логгер с уровнями и фильтрацией
    proxy/
        server.go               # Основной сервер
        listener.go             # Логирующий listener и conn
 .env                   # Пример конфигурации
 Dockerfile                       # Docker образ
 Makefile                         # Сборка и установка
 go.mod                           # Go модуль
 README.md                        # Эта документация
```

---

##  Разработка

### Команды Makefile

```
make help         # Показать все команды
make deps         # Установить зависимости
make build        # Собрать бинарник
make run          # Собрать и запустить
make run-debug    # Собрать и запустить с debug
make run-dev      # Запуск в dev режиме (полное логирование)
make clean        # Очистить собранные файлы
make test         # Запустить тесты
make install      # Установить в систему как сервис
make uninstall    # Удалить из системы
make setup         # Сгенерировать .env файл с случайными учетными данными
make generate-env  # Только генерация .env файла
```

### Тестирование

```bash
# Запустить все тесты
make test

# Или напрямую через go
go test ./...

# С подробным выводом
go test -v ./...

# Только конкретный пакет
go test -v ./internal/config/
go test -v ./internal/logger/
go test -v ./internal/proxy/

# Один конкретный тест
go test -v -run TestConfigValidate ./internal/config/
go test -v -run TestServerListens ./internal/proxy/

# С отчётом о покрытии
go test -cover ./...
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
```

### Сборка вручную

```
# Инициализация модуля
go mod init go-socks5-relay

# Установка зависимостей
go mod tidy

# Сборка
go build -o build/socks5-proxy ./cmd/socks5-proxy

# Запуск
./build/socks5-proxy -debug
```

---

##  Устранение проблем

### Проблема: нет логов подключений

**Решение:** Запустите с флагом `-debug` или `-log-level debug`

```
./build/socks5-proxy -debug
```

### Проблема: "Connection refused"

**Проверьте:**

```
# Слушает ли сервер
sudo netstat -tlnp | grep 5431

# Не блокирует ли фаервол
sudo ufw status
```

### Проблема: "authentication failed"

**Проверьте** логин и пароль в `.env`:

```
cat .env
```

### Проблема: порт уже используется

```
# Найти процесс, использующий порт
sudo lsof -i :1080

# Остановить процесс или изменить порт
./build/socks5-proxy -port 1080
```

### Проблема: ошибка сборки "cannot find package"

```
# Очистить кэш модулей
go clean -modcache

# Повторить установку зависимостей
go mod tidy
```

---

##  Для LLM (подсказка)

Этот README можно использовать как промт для LLM. Если вам нужно обратиться за помощью по этому проекту, используйте следующий формат:

---
У меня есть проект SOCKS5 прокси на Go со следующей структурой:
- cmd/socks5-proxy/main.go  точка входа, флаги командной строки
- internal/config/  конфигурация из файла .env
- internal/logger/  логгер с уровнями (error/warn/info/debug)
- internal/proxy/  сервер и логирующий listener
- Makefile  сборка и установка
- Dockerfile  контейнеризация
- systemd сервис  socks5-proxy.service

Что мне нужно сделать: [ваш вопрос]
---

---

##  Лицензия

---

**Вопросы или предложения?** Создайте [Issue](https://github.com/yourusername/go-socks5-relay/issues) на GitHub.
