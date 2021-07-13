package log

type Logger interface {
	// Info uses fmt.Sprint to construct and log a message.
	Info(args ...interface{})

	// Infof uses fmt.Sprintf to log a template message.
	Infof(msgFormat string, args ...interface{})

	// Infow uses fmt.Sprintf to log a template message with extra context value.
	Infow(keysAndValues []interface{}, msgFormat string, args ...interface{})

	// Debug uses fmt.Sprint to construct and log a message.
	Debug(args ...interface{})

	// Debugf uses fmt.Sprintf to log a template message.
	Debugf(msgFormat string, args ...interface{})

	// Debugw uses fmt.Sprintf to log a template message with extra context value.
	Debugw(keysAndValues []interface{}, msgFormat string, args ...interface{})

	// Warn uses fmt.Sprint to construct and log a message.
	Warn(args ...interface{})

	// Warnf uses fmt.Sprintf to log a template message.
	Warnf(msgFormat string, args ...interface{})

	// Warnw uses fmt.Sprintf to log a template message with extra context value.
	Warnw(keysAndValues []interface{}, msgFormat string, args ...interface{})

	// Error uses fmt.Sprint to construct and log a message.
	Error(args ...interface{})

	// Errorf uses fmt.Sprintf to log a template message.
	Errorf(msgFormat string, args ...interface{})

	// Errorw uses fmt.Sprintf to log a template message with extra context value.
	Errorw(keysAndValues []interface{}, msgFormat string, args ...interface{})

	// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
	Fatal(args ...interface{})

	// Fatalf uses fmt.Sprintf to log a template message, then calls os.Exit.
	Fatalf(msgFormat string, args ...interface{})

	// Fatalw uses fmt.Sprintf to log a template message with extra context value, then calls os.Exit.
	Fatalw(keysAndValues []interface{}, msgFormat string, args ...interface{})
}
