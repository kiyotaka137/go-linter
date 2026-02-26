package badlogs_logslog

import (
	log "github.com/sirupsen/logrus"
	"log/slog"
)

func run(password string) {
	log.Info("Starting server on port 8080")      // want "log message must start with lowercase letter" "log message have forbidden symbols"
	log.Error("ошибка подключения к базе данных") // want "log message have forbidden symbols"
	log.Warn("connection failed!!!")              // want "log message have forbidden symbols"
	log.Debug("api_key leaked")                   // want "log message have sensitive keyword \"api_key\"" "log message have forbidden symbols"
	log.Info("user card 1234567812345678")        // want "log message matches sensitive pattern"

	slog.Info("Server started!")                  // want "log message must start with lowercase letter" "log message have forbidden symbols"
	slog.Error("пароль пользователя невалиден")   // want "log message have forbidden symbols"
	slog.Warn("warning: something went wrong...") // want "log message have forbidden symbols"
	slog.Debug("token exposed")                   // want "log message have sensitive keyword \"token\""

	_ = password
}
