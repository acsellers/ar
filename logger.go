package db

// I need to write the code for this, while I'm at I can
// add a quick performance tweak of a prepared query cache
type Logger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}
