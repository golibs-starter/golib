package golib

import (
	"fmt"
	"github.com/golibs-starter/golib/log"
	"github.com/golibs-starter/golib/log/field"
	"go.uber.org/fx/fxevent"
	"strings"
)

// FxLogger is a Fx event logger that logs events to field.
type FxLogger struct {
	logger log.Logger
}

func NewFxLogger(logger log.Logger) fxevent.Logger {
	return &FxLogger{logger: logger.WithField(field.String("module", "fx"))}
}

func (l *FxLogger) logf(msg string, args ...interface{}) {
	l.logger.Debugf(msg, args...)
}

func (l *FxLogger) errorf(msg string, args ...interface{}) {
	l.logger.Errorf(msg, args...)
}

// LogEvent logs the given event to the provided Zap logger.
func (l *FxLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.logf("HOOK OnStart: %s executing (caller: %s)", e.FunctionName, e.CallerName)
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.logf("HOOK OnStart: %s called by %s failed in %s: %+v", e.FunctionName, e.CallerName, e.Runtime, e.Err)
		} else {
			l.logf("HOOK OnStart: %s called by %s ran successfully in %s", e.FunctionName, e.CallerName, e.Runtime)
		}
	case *fxevent.OnStopExecuting:
		l.logf("HOOK OnStop: %s executing (caller: %s)", e.FunctionName, e.CallerName)
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.logf("HOOK OnStop: %s called by %s failed in %s: %+v", e.FunctionName, e.CallerName, e.Runtime, e.Err)
		} else {
			l.logf("HOOK OnStop: %s called by %s ran successfully in %s", e.FunctionName, e.CallerName, e.Runtime)
		}
	case *fxevent.Supplied:
		if e.Err != nil {
			l.errorf("ERROR: Failed to supply %v: %+v", e.TypeName, e.Err)
		} else if e.ModuleName != "" {
			l.logf("SUPPLY: %v from module %q", e.TypeName, e.ModuleName)
		} else {
			l.logf("SUPPLY: %v", e.TypeName)
		}
	case *fxevent.Provided:
		var privateStr string
		if e.Private {
			privateStr = " (PRIVATE)"
		}
		for _, rtype := range e.OutputTypeNames {
			if e.ModuleName != "" {
				l.logf("PROVIDE%v: %v <= %v from module %q", privateStr, rtype, e.ConstructorName, e.ModuleName)
			} else {
				l.logf("PROVIDE%v: %v <= %v", privateStr, rtype, e.ConstructorName)
			}
		}
		if e.Err != nil {
			l.logf("Error after options were applied: %+v", e.Err)
		}
	case *fxevent.Replaced:
		for _, rtype := range e.OutputTypeNames {
			if e.ModuleName != "" {
				l.logf("REPLACE: %v from module %q", rtype, e.ModuleName)
			} else {
				l.logf("REPLACE: %v", rtype)
			}
		}
		if e.Err != nil {
			l.errorf("Failed to replace: %+v", e.Err)
		}
	case *fxevent.Decorated:
		for _, rtype := range e.OutputTypeNames {
			if e.ModuleName != "" {
				l.logf("DECORATE: %v <= %v from module %q", rtype, e.DecoratorName, e.ModuleName)
			} else {
				l.logf("DECORATE: %v <= %v", rtype, e.DecoratorName)
			}
		}
		if e.Err != nil {
			l.logf("Error after options were applied: %+v", e.Err)
		}
	case *fxevent.Run:
		var moduleStr string
		if e.ModuleName != "" {
			moduleStr = fmt.Sprintf(" from module %q", e.ModuleName)
		}
		l.logf("RUN: %v: %v%v", e.Kind, e.Name, moduleStr)
		if e.Err != nil {
			l.errorf("Error returned: %+v", e.Err)
		}

	case *fxevent.Invoking:
		if e.ModuleName != "" {
			l.logf("INVOKE: %s from module %q", e.FunctionName, e.ModuleName)
		} else {
			l.logf("INVOKE: %s", e.FunctionName)
		}
	case *fxevent.Invoked:
		if e.Err != nil {
			l.errorf("fx.Invoke(%v) called from: %+vFailed: %+v", e.FunctionName, e.Trace, e.Err)
		}
	case *fxevent.Stopping:
		l.logf("%v", strings.ToUpper(e.Signal.String()))
	case *fxevent.Stopped:
		if e.Err != nil {
			l.errorf("Failed to stop cleanly: %+v", e.Err)
		}
	case *fxevent.RollingBack:
		l.errorf("Start failed, rolling back: %+v", e.StartErr)
	case *fxevent.RolledBack:
		if e.Err != nil {
			l.errorf("Couldn't roll back cleanly: %+v", e.Err)
		}
	case *fxevent.Started:
		if e.Err != nil {
			l.errorf("Failed to start: %+v", e.Err)
		} else {
			l.logf("RUNNING")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.errorf("Failed to initialize custom logger: %+v", e.Err)
		} else {
			l.logf("LOGGER: Initialized custom logger from %v", e.ConstructorName)
		}
	}
}
