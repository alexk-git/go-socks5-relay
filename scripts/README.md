# scripts/

Вспомогательные shell-скрипты для настройки и запуска SOCKS5 прокси.

## Содержание

| Файл | Назначение |
|---|---|
| [`generate-env.sh`](#generate-envsh) | Генерация `config/.env` со случайными учётными данными |
| [`entrypoint.sh`](#entrypointsh) | Точка входа Docker-контейнера |
| [`generate-compose-env.sh`](#generate-compose-envsh) | Генерация корневого `.env` для docker-compose |

---

### `generate-env.sh`

**Назначение:** Создаёт файл `config/.env` на основе шаблона `config/.env.example`, заполняя случайные значения `username` и `password`.

**Логика работы:**
1. Проверяет, существует ли `config/.env` — если да, пропускает генерацию.
2. Копирует `config/.env.example` в `config/.env`.
3. Генерирует случайный username в формате `user_<8hex>` и пароль длиной 24 hex-символа.
4. Подставляет их в `.env` через `sed` (совместимо с GNU sed и BSD sed/macOS).

**Используется в:**
- `Makefile` — цели `generate-env`, `setup`, `run`, `run-debug`, `run-dev`, `docker-up`
- `entrypoint.sh` — при старте контейнера, если `.env` отсутствует
- `generate-compose-env.sh` — если `config/.env` не найден

**Зависимости:** `openssl` (или `/dev/urandom` как fallback), `sed`

**Пример:**
```bash
# Стандартный запуск
./scripts/generate-env.sh

# С указанием путей вручную
./scripts/generate-env.sh config/.env.example config/.env
```

---

### `entrypoint.sh`

**Назначение:** Точка входа Docker-контейнера. Гарантирует, что перед запуском прокси существует файл конфигурации.

**Логика работы:**
1. Проверяет наличие `/app/config/.env`.
2. Если файла нет — запускает `generate-env.sh`, используя шаблон из `/app/config/.env.example` (или `/app/.env.example` как fallback).
3. Запускает бинарник `socks5-proxy` через `exec`, передавая ему все аргументы командной строки (`$@`).

**Используется в:**
- `Dockerfile` — установлен как `ENTRYPOINT`, выполняется при каждом старте контейнера

**Пример (в Dockerfile):**
```dockerfile
ENTRYPOINT ["/app/scripts/entrypoint.sh"]
CMD ["-log-level", "debug"]
```

---

### `generate-compose-env.sh`

**Назначение:** Создаёт корневой файл `.env` с переменной `SOCKS5_PORT` для использования в `docker-compose.yml`.

**Логика работы:**
1. Если `config/.env` отсутствует — вызывает `generate-env.sh` для его создания.
2. Читает порт из `config/.env` (или `config/.env.example` как fallback).
3. Создаёт в корне проекта файл `.env` с содержимым:
   ```
   SOCKS5_PORT=3128
   ```

**Используется в:**
- `Makefile` — цель `generate-compose-env`

**Пример:**
```bash
./scripts/generate-compose-env.sh
docker-compose up
```
