package transfer

// Logger is used to log critical error messages.
type Logger interface {
	Print(v ...interface{})
}
