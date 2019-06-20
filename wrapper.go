package unilogger

import (
	"errors"
	"fmt"
)

// LoggerWrapper is wrapper around any third-party logger that implements any of
// the following interfaces:
//	# ExtendedLeveledLogger
//	# ShortLeveledLogger
//	# LeveledLogger
//	# StdLogger
// By wrapping the logger it implements ExtendedLeveledLogger.
// For loggers that implements only StdLogger, LoggerWrapper tries to virtualize
// leveled logger behaviour. It simply adds level name before logging message.
// If a logger implements LeveledLogger that doesn't have specific log line '****ln()' methods,
// it uses default non 'ln' functions - i.e. instead 'Infoln' uses 'Info'.
type LoggerWrapper struct {
	logger        interface{}
	currentLogger int
}

// NewLoggerWrapper creates a LoggerWrapper wrapper over provided 'logger' argument
// By default the function checks if provided logger implements logging interfaces
// in a following hierarchy:
//	# ExtendedLeveledLogger
//	# ShortLeveledLogger
//	# LeveledLogger
//	# StdLogger
// if logger doesn't implement an interface it tries to check the next in hierarchy.
// If it doesn't implement any of known logging interfaces the function returns error.
func NewLoggerWrapper(logger interface{}) (*LoggerWrapper, error) {
	return newLoggerWrapper(logger)
}

// MustGetLoggerWrapper creates a LoggerWrapper wrapper over provided 'logger' argument.
// By default the function checks if provided logger implements logging interfaces
// in a following hierarchy:
//	# ExtendedLeveledLogger
//	# ShortLeveledLogger
//	# LeveledLogger
//	# StdLogger
// if logger doesn't implement an interface it tries to check the next in hierarchy.
// If it doesn't implement any of known logging interfaces the function panics.
func MustGetLoggerWrapper(logger interface{}) *LoggerWrapper {
	wrapper, err := newLoggerWrapper(logger)
	if err != nil {
		panic(err)
	}
	return wrapper
}

func newLoggerWrapper(logger interface{}) (*LoggerWrapper, error) {
	wrapper := &LoggerWrapper{}
	var err error

	if l, ok := logger.(ExtendedLeveledLogger); ok {
		wrapper.logger = l
		wrapper.currentLogger = 4
		return wrapper, nil
	}

	if l, ok := logger.(ShortLeveledLogger); ok {
		wrapper.logger = l
		wrapper.currentLogger = 3
		return wrapper, nil
	}

	if l, ok := logger.(LeveledLogger); ok {
		wrapper.logger = l
		wrapper.currentLogger = 2
		return wrapper, nil
	}

	if l, ok := logger.(StdLogger); ok {
		wrapper.logger = l
		wrapper.currentLogger = 1
		return wrapper, nil
	}

	err = errors.New("Provided logger doesn't implement any known interfaces")
	return nil, err
}

// Print logs a message.
// Arguments are handled in the manner of log.Print for StdLogger and
// Extended LeveledLogger as well as log.Info for LeveledLogger
func (c *LoggerWrapper) Print(args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		log.Print(args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Info(args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Info(args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Print(args...)
	default:
	}
}

// Printf logs a formatted message.
// Arguments are handled in the manner of log.Printf for StdLogger and
// Extended LeveledLogger as well as log.Infof for LeveledLogger
func (c *LoggerWrapper) Printf(format string, args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		log.Printf(format, args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Infof(format, args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Infof(format, args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Printf(format, args...)
	default:
	}
}

// Println logs a message.
// Arguments are handled in the manner of log.Println for StdLogger and
// Extended LeveledLogger as well as log.Info for LeveledLogger
func (c *LoggerWrapper) Println(args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		log.Println(args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Info(args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Info(args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Println(args...)
	default:

	}
}

// Debug logs a message with DEBUG level.
// Arguments are handled in the manner of log.Print for StdLogger,
// log.Debug for ExtendedLeveledLogger and LeveledLogger.
func (c *LoggerWrapper) Debug(args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(DEBUG, nil, args...)
		log.Print(args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Debug(args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Debug(args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Debug(args...)
	default:
	}
}

// Debugf logs a formatted message with DEBUG level.
// Arguments are handled in the manner of log.Printf for StdLogger,
// log.Debugf for ExtendedLeveledLogger, ShortLeveledLogger and LeveledLogger.
func (c *LoggerWrapper) Debugf(format string, args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(DEBUG, &format, args...)
		log.Printf(format, args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Debugf(format, args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Debugf(format, args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Debugf(format, args...)
	default:
	}
}

// Debugln logs a message with DEBUG level.
// Arguments are handled in the manner of log.Println for StdLogger,
// log.Debugln for ExtendedLeveledLogger and log.Debug for LeveledLogger and ShortLeveledLogger.
func (c *LoggerWrapper) Debugln(args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(DEBUG, nil, args...)
		log.Println(args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Debug(args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Debug(args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Debugln(args...)
	default:
	}

}

// Info logs a message with INFO level.
// Arguments are handled in the manner of log.Print for StdLogger,
// log.Info for ExtendedLeveledLogger, ShortLeveledLogger and LeveledLogger.
func (c *LoggerWrapper) Info(args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(INFO, nil, args...)
		log.Print(args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Info(args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Info(args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Info(args...)
	default:
	}

}

// Infof logs a formatted message with INFO level.
// Arguments are handled in the manner of log.Printf for StdLogger,
// log.Infof for ExtendedLeveledLogger, ShortLeveledLogger and LeveledLogger.
func (c *LoggerWrapper) Infof(format string, args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(INFO, &format, args...)
		log.Printf(format, args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Infof(format, args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Infof(format, args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Infof(format, args...)
	default:
	}
}

// Infoln logs a message with INFO level.
// Arguments are handled in the manner of log.Println for StdLogger,
// log.Infoln for ExtendedLeveledLogger and log.Info for LeveledLogger and ShortLeveledLogger.
func (c *LoggerWrapper) Infoln(args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(INFO, nil, args...)
		log.Println(args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Info(args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Info(args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Infoln(args...)
	default:
	}
}

// Warning logs a message with WARNING level.
// Arguments are handled in the manner of log.Print for StdLogger,
// log.Warning for ExtendedLeveledLogger, LeveledLogger and
// log.Warn for ShortLeveledLogger.
func (c *LoggerWrapper) Warning(args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(WARNING, nil, args...)
		log.Print(args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Warning(args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Warn(args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Warning(args...)
	default:
	}
}

// Warningf logs a formatted message with WARNING level.
// Arguments are handled in the manner of log.Printf for StdLogger,
// log.Warningf for ExtendedLeveledLogger, LeveledLogger and log.Warnf for ShortLeveledLogger.
func (c *LoggerWrapper) Warningf(format string, args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(WARNING, &format, args...)
		log.Printf(format, args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Warningf(format, args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Warnf(format, args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Warningf(format, args...)
	default:
	}
}

// Warningln logs a message with WARNING level.
// Arguments are handled in the manner of log.Println for StdLogger,
// log.Warningln for ExtendedLeveledLogger, log.Warning for LeveledLogger
// and log.Warn for ShortLeveledLogger.
func (c *LoggerWrapper) Warningln(args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(WARNING, nil, args...)
		log.Println(args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Warning(args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Warn(args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Warningln(args...)
	default:
	}
}

// Error logs a message with ERROR level.
// Arguments are handled in the manner of log.Print for StdLogger,
// log.Error for ExtendedLeveledLogger, LeveledLogger and ShortLeveledLogger.
func (c *LoggerWrapper) Error(args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(ERROR, nil, args...)
		log.Print(args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Error(args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Error(args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Error(args...)
	default:
	}
}

// Errorf logs a formatted message with ERROR level.
// Arguments are handled in the manner of log.Printf for StdLogger,
// log.Errorf for ExtendedLeveledLogger, LeveledLogger and ShortLeveledLogger.
func (c *LoggerWrapper) Errorf(format string, args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(ERROR, &format, args...)
		log.Printf(format, args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Errorf(format, args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Errorf(format, args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Errorf(format, args...)
	default:
	}
}

// Errorln logs a message with ERROR level.
// Arguments are handled in the manner of log.Println for StdLogger,
// log.Debugln for ExtendedLeveledLogger and log.Error for LeveledLogger and ShortLeveledLogger.
func (c *LoggerWrapper) Errorln(args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(ERROR, nil, args...)
		log.Println(args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Error(args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Error(args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Errorln(args...)
	default:
	}
}

// Fatal logs a message with CRITICAL level. Afterwards it should excute os.Exit(1).
// Arguments are handled in the manner of log.Fatal for StdLogger, LeveledLogger,
// ShortLeveledLogger and ExtendedLeveledLogger.
func (c *LoggerWrapper) Fatal(args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(CRITICAL, nil, args...)
		log.Fatal(args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Fatal(args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Fatal(args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Fatal(args...)
	default:
	}
}

// Fatalf logs a formatted message with CRITICAL level. Afterwards it should excute os.Exit(1).
// Arguments are handled in the manner of log.Fatalf for StdLogger, LeveledLogger,
// ShortLeveledLogger and ExtendedLeveledLogger.
func (c *LoggerWrapper) Fatalf(format string, args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(CRITICAL, &format, args...)
		log.Fatalf(format, args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Fatalf(format, args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Fatalf(format, args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Fatalf(format, args...)
	default:
	}
}

// Fatalln logs a message with CRITICAL level. Afterwards it should excute os.Exit(1).
// Arguments are handled in the manner of log.Fatalln for StdLogger and ExtendedLeveldLogger,
// and log.Fatal for LeveledLogger and ShortLeveledLogger.
func (c *LoggerWrapper) Fatalln(args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(CRITICAL, nil, args...)
		log.Fatalln(args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Fatal(args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Fatal(args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Fatalln(args...)
	default:
	}
}

// Panic logs a message with CRITICAL level. Afterwards it should panic.
// Arguments are handled in the manner of log.Panic for StdLogger, LeveledLogger,
// ShortLeveledLogger and ExtendedLeveledLogger .
func (c *LoggerWrapper) Panic(args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(CRITICAL, nil, args...)
		log.Panic(args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Panic(args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Panic(args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Panic(args...)
	default:
	}
}

// Panicf logs a formatted message with CRITICAL level. Afterwards it should panic.
// Arguments are handled in the manner of log.Panicf for StdLogger, LeveledLogger,
// ShortLeveledLogger and ExtendedLeveledLogger.
func (c *LoggerWrapper) Panicf(format string, args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(CRITICAL, &format, args...)
		log.Panicf(format, args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Panicf(format, args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Panicf(format, args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Panicf(format, args...)
	default:
	}
}

// Panicln logs a message with CRITICAL level. Afterwards it should panic.
// Arguments are handled in the manner of log.Panicln for StdLogger and ExtendedLeveledLogger,
// and log.Panic LeveledLogger and ShortLeveledLogger.
func (c *LoggerWrapper) Panicln(args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(CRITICAL, nil, args...)
		log.Panicln(args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Panic(args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Panic(args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Panicln(args...)
	default:
	}
}

func buildLeveled(level Level, format *string, args ...interface{}) (leveled []interface{}) {
	if format == nil {
		leveled = append(leveled, fmt.Sprintf("%s: ", level))
		leveled = append(leveled, args...)
	} else {
		leveled = append(leveled, args...)
		msg := fmt.Sprintf("%s: %s", level, *format)
		*format = msg
	}
	return leveled
}
