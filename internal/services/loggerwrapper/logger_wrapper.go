package loggerwrapper

type Logger interface {
	Debug(msg string, values ...interface{})
	Info(msg string, values ...interface{})
	Warn(msg string, values ...interface{})
	Error(msg string, values ...interface{})
	Fatal(msg string, values ...interface{})
	Level() string
}
