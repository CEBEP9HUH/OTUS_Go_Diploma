package loggerwrapper

type emptyLogger struct{}

func NewEmptyLogger() Logger {
	return emptyLogger{}
}

func (l emptyLogger) Debug(string, ...interface{}) {
}

func (l emptyLogger) Info(string, ...interface{}) {
}

func (l emptyLogger) Warn(string, ...interface{}) {
}

func (l emptyLogger) Error(string, ...interface{}) {
}

func (l emptyLogger) Fatal(string, ...interface{}) {
}

func (l emptyLogger) Level() string {
	return ""
}
