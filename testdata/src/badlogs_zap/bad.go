package badlogs_zap

import zaplib "go.uber.org/zap"

var zap = zaplib.NewExample().Sugar()

func run() {
	zap.Info("Failed to connect to database") // want "log message must start with lowercase letter" "log message have forbidden symbols"
	zap.Error("запуск сервера")               // want "log message have forbidden symbols"
	zap.Warn("server started! 🚀")             // want "log message have forbidden symbols"
	zap.Debug("client_secret exposed")        // want "log message have sensitive keyword \"secret\"" "log message have forbidden symbols"
}
