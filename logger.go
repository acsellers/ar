package db

type Logger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}
