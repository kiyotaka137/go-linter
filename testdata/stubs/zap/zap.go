package zap

type Logger struct{}
type SugaredLogger struct{}

func NewExample(...Option) *Logger {
	return &Logger{}
}

func (l *Logger) Sugar() *SugaredLogger {
	return &SugaredLogger{}
}

func (l *SugaredLogger) Info(string, ...any)  {}
func (l *SugaredLogger) Warn(string, ...any)  {}
func (l *SugaredLogger) Error(string, ...any) {}
func (l *SugaredLogger) Debug(string, ...any) {}

type Option interface{}
