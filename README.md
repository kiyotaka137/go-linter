# loglint (Go analyzer)

Линтер для проверки лог-сообщений в Go-коде на основе `golang.org/x/tools/go/analysis`.

## Что сделано по

### Этап 1. Базовая структура
- Создан проект на Go с анализатором `Analyzer` (`analyzer` пакет).
- Реализован AST-обход и поиск лог-вызовов.
- Подключены тесты через `analysistest`.

### Этап 2. Реализация правил
Сейчас проверяются:
- Сообщение начинается со строчной буквы.
- Сообщение не содержит спецсимволы/эмодзи.
- Сообщение не содержит чувствительные данные по ключевым словам.
- Сообщение не содержит чувствительные данные по regex-паттернам из конфига.

Примечание про правило "только английский":
- Отдельного конкретного правила `english` нет.
- На практике это покрывается правилом `symbols`: разрешены только `[a-z0-9]`и пробел, поэтому кириллица и прочий не ASCII текст отклоняются.

### Этап 3. Тестирование
- Добавлены unit-тесты для `internal/rules`.
- Добавлены тестовые пакеты в `testdata/src/...` для `analysistest`.
- Проверка анализатора выполняется в `analyzer/analyzer_test.go`.

### Этап 4. Интеграция
- Анализатор реализован в формате `go/analysis` (база для интеграции в экосистему линтеров).
- Подготовлен базовый CI в GitHub Actions: [ci.yml](.github/workflows/ci.yml).

## Бонусы
- Конфиг правил через YAML (`config.yaml`):
  - включение/выключение правил;
  - ключевые слова sensitive;
  - regex-паттерны sensitive.
- SuggestedFixes:
  - автоисправление для `lowercase`;
  - автоисправление для `symbols`.
- Кастомные паттерны sensitive через `rules.sensitive.patterns`.

## Конфигурация

Пример `config.yaml`:

```yaml
rules:
  lowercase:
    enabled: true
  symbols:
    enabled: true
  sensitive:
    enabled: true
    keywords:
      - password
      - token
    patterns:
      - '\b\d{16}\b'
```

Также поддерживается флаг анализатора:
- `-config=/path/to/config.yaml`

## Интеграция с golangci-lint

Ниже пример того, что нужно для запуска линтера как custom plugin в `golangci-lint` (Module Plugin System).
В проекте уже есть plugin entrypoint: [plugin.go](plugin.go) (регистрация линтера `loglint`).

### 1. Как собрать/подключить плагин

Создать `.custom-gcl.yml` в корне проекта:

```yaml
version: v2.8.0
name: custom-gcl
plugins:
  - module: github.com/your-org/go-linter # замените на ваш module path из go.mod
    path: .
```

Собрать кастомный бинарник:

```bash
golangci-lint custom
```

После этого появится бинарник `./custom-gcl`.

### 2. Как включить линтер в конфиге golangci-lint

Пример `.golangci.yml`:

```yaml
version: "2"

linters:
  default: none
  enable:
    - loglint
  settings:
    custom:
      loglint:
        type: module
        description: Checks log messages policy
        settings:
          config: ./config.yaml
```

### 3. Пример команды запуска

```bash
./custom-gcl run ./...
```

### 4. Пример `config.yaml` для правил

```yaml
rules:
  lowercase:
    enabled: true
  symbols:
    enabled: true
  sensitive:
    enabled: true
    keywords:
      - password
      - token
      - client_secret
    patterns:
      - '\b\d{16}\b'
      - '(?i)bearer\s+[a-z0-9._-]+'
```

### 5. Типичный вывод

```text
testdata/src/badlogs_logslog/bad.go:9:11: loglint: log message must start with lowercase letter
testdata/src/badlogs_logslog/bad.go:12:11: loglint: log message have sensitive keyword "api_key"
testdata/src/badlogs_logslog/bad.go:13:11: loglint: log message matches sensitive pattern "\b\d{16}\b"
testdata/src/badlogs_zap/bad.go:10:10: loglint: log message have forbidden symbols
```

## Поддерживаемые логгеры
- `log/slog`
- `go.uber.org/zap` (вызовы через идентификатор `zap`)

## Как запустить тесты

```bash
go test ./...
```

Отдельно анализатор:

```bash
go test ./analyzer -v
```

## Текущие ограничения
- Проверяются только сообщения, которые удаётся статически извлечь как строку.
- В основном покрыты строковые литералы и конкатенация строковых литералов.
- Динамические выражения вида `"token: " + tokenVar` сейчас не анализируются полностью.
