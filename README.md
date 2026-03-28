# SOCKS5 Proxy Server on Go

[![Go Version](https://img.shields.io/badge/Go-1.22+-blue.svg)](https://golang.org)
[![Systemd](https://img.shields.io/badge/systemd-ready-green.svg)](https://systemd.io)

Производственно-готовый SOCKS5 прокси-сервер с аутентификацией, написанный на Go. Поддерживает graceful shutdown, гибкую настройку через флаги и переменные окружения, а также установку как systemd сервис.

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

---

##  Настройка конфигурации

### Файл конфигурации (`env.properties`)

Создайте файл `env.properties` в корневой директории проекта со следующим содержимым:

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

### Приоритет настроек

Настройки применяются в следующем порядке (от высшего приоритета к низшему):

| Приоритет | Источник | Пример |
|-----------|----------|--------|
| 1 | Флаги командной строки | `-port 8080` |
| 2 | Переменные окружения | `SOCKS5_CONFIG=/path/to/config` |
| 3 | Файл конфигурации | `env.properties` |
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
export SOCKS5_CONFIG=/etc/my-proxy/env.properties

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

Перед первой сборкой необходимо инициализировать модуль и скачать зависимости:

```
# Инициализация модуля (если go.mod отсутствует)
go mod init go-socks5-relay

# Скачивание зависимостей
go mod tidy
```

Или используйте Makefile:

```
make deps
```

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
[socks5] [I]: Конфигурация: env.properties
[socks5] [I]: Уровень логирования: info
[socks5] [I]: Debug режим: false
[socks5] [I]: ===========================
[socks5] [I]: SOCKS5-прокси запущен на 0.0.0.0:1080 (TCP + UDP)
```

В другом терминале проверьте работу прокси:

```
# Замените username:password на ваши данные из env.properties
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
./build/socks5-proxy -config /etc/my-proxy/env.properties

# Комбинация
./build/socks5-proxy -debug -port 1080 -log-level debug
```

### Переменные окружения

```
# Включить debug режим
export SOCKS5_DEBUG=1

# Указать путь к конфигурации
export SOCKS5_CONFIG=/etc/proxy/env.properties

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
sudo cp env.properties /etc/socks5-proxy/
sudo chmod 600 /etc/socks5-proxy/env.properties

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
ExecStart=/usr/local/bin/socks5-proxy -config /etc/socks5-proxy/env.properties
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
docker run -p 1080:1080 -v $(pwd)/.env:/app/env.properties socks5-proxy

# В фоновом режиме
docker run -d -p 1080:1080 --name socks5-proxy socks5-proxy -debug
```

### Проверка

```
curl --proxy socks5://username:password@localhost:1080 https://api.ipify.org
```

---

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
 env.properties                   # Пример конфигурации
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

**Проверьте** логин и пароль в `env.properties`:

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
- internal/config/  конфигурация из файла env.properties
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
