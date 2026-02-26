module go-linter-testdata

go 1.25.7

require (
	github.com/sirupsen/logrus v0.0.0
	go.uber.org/zap v0.0.0
)

replace github.com/sirupsen/logrus => ./stubs/logrus

replace go.uber.org/zap => ./stubs/zap
