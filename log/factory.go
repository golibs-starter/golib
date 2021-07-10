package log

var globalLogger *logger

// RegisterGlobal Register a logger instance as global
func RegisterGlobal(logger *logger) {
	globalLogger = logger
}

// GetGlobalLogger Get global logger instance
func GetGlobalLogger() *logger {
	return globalLogger
}

// Info uses fmt.Sprint to construct and log a message.
func Info(args ...interface{}) {
	globalLogger.Info(args...)
}

// Infof uses fmt.Sprintf to log a template message.
func Infof(msgFormat string, args ...interface{}) {
	globalLogger.Infof(msgFormat, args...)
}

// Infow uses fmt.Sprintf to log a template message with extra context value.
func Infow(keysAndValues []interface{}, msgFormat string, args ...interface{}) {
	globalLogger.Infow(keysAndValues, msgFormat, args...)
}

// Debug uses fmt.Sprint to construct and log a message.
func Debug(args ...interface{}) {
	globalLogger.Debug(args...)
}

// Debugf uses fmt.Sprintf to log a template message.
func Debugf(msgFormat string, args ...interface{}) {
	globalLogger.Debugf(msgFormat, args...)
}

// Debugw uses fmt.Sprintf to log a template message with extra context value.
func Debugw(keysAndValues []interface{}, msgFormat string, args ...interface{}) {
	globalLogger.Debugw(keysAndValues, msgFormat, args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func Warn(args ...interface{}) {
	globalLogger.Warn(args...)
}

// Warnf uses fmt.Sprintf to log a template message.
func Warnf(msgFormat string, args ...interface{}) {
	globalLogger.Warnf(msgFormat, args...)
}

// Warnw uses fmt.Sprintf to log a template message with extra context value.
func Warnw(keysAndValues []interface{}, msgFormat string, args ...interface{}) {
	globalLogger.Warnw(keysAndValues, msgFormat, args...)
}

// Error uses fmt.Sprint to construct and log a message.
func Error(args ...interface{}) {
	globalLogger.Error(args...)
}

// Errorf uses fmt.Sprintf to log a template message.
func Errorf(msgFormat string, args ...interface{}) {
	globalLogger.Errorf(msgFormat, args...)
}

// Errorw uses fmt.Sprintf to log a template message with extra context value.
func Errorw(keysAndValues []interface{}, msgFormat string, args ...interface{}) {
	globalLogger.Errorw(keysAndValues, msgFormat, args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(args ...interface{}) {
	globalLogger.Fatal(args...)
}

// Fatalf uses fmt.Sprintf to log a template message, then calls os.Exit.
func Fatalf(msgFormat string, args ...interface{}) {
	globalLogger.Fatalf(msgFormat, args...)
}

// Fatalw uses fmt.Sprintf to log a template message with extra context value, then calls os.Exit.
func Fatalw(keysAndValues []interface{}, msgFormat string, args ...interface{}) {
	globalLogger.Fatalw(keysAndValues, msgFormat, args...)
}
