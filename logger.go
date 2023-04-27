// Package uzap contains basic configuration for go.uber.org/zap, for easy use with 12-factor apps.
package uzap

import (
	"os"

	"github.com/Unity-Technologies/zapdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/kelseyhightower/envconfig"
)

// Options is used to parse environment vars with the log level and optional debug flag.
type Options struct {
	Level zapcore.Level // zap defaults to INFO
	Debug bool          // defaults to false, displays human readable output instead of json
}

// NewZapLogger configures a zap.Logger for use in container based environments
// ERROR level logs are written to stderr and all other levels are written to stdout
// Useful in Kubernetes where stderr & stdout are interpreted as ERROR & INFO level logs respectively
// opts.Debug controls the loggers output. Human readable when true; JSON when false.
func NewZapLogger(opts *Options) (*zap.Logger, zap.AtomicLevel) {
	if opts == nil {
		opts = &Options{}
	}

	level := zap.NewAtomicLevelAt(opts.Level)

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	// It is useful for Kubernetes deployment.
	// Kubernetes interprets os.Stdout log items as INFO and os.Stderr log items
	// as ERROR by default.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= level.Level() && lvl < zapcore.ErrorLevel
	})

	// Output channels
	consoleInfos := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	// Setup Config
	var (
		ecfg zapcore.EncoderConfig
		enc  zapcore.Encoder
	)

	if opts.Debug {
		ecfg = zapdriver.NewDevelopmentEncoderConfig()
		enc = zapcore.NewConsoleEncoder(ecfg)
	} else {
		ecfg = zapdriver.NewProductionEncoderConfig()
		enc = zapcore.NewJSONEncoder(ecfg)
	}

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.
	core := zapcore.NewTee(
		zapcore.NewCore(enc, consoleErrors, highPriority),
		zapcore.NewCore(enc, consoleInfos, lowPriority),
	)
	// From a zapcore.Core, it's easy to construct a Logger.
	return zap.New(core), level
}

// MustZap is an ease of use function that replaces zap globals
// and redirects standard `package log` output to a new zap logger.
// It returns a deferrable function, for calling zap.Logger.Sync at program termination.
func MustZap() func() {
	return MustZapWithLevel(zapcore.InfoLevel)
}

// MustZapWithLevel is an ease of use function that replaces zap globals
// and redirects standard `package log` output to a new zap logger.
// It returns a deferrable function, for calling zap.Logger.Sync at program termination.
func MustZapWithLevel(lvl zapcore.Level) func() {
	opts := &Options{Level: lvl}
	if err := envconfig.Process("log", opts); err != nil {
		panic(err)
	}

	logger, _ := NewZapLogger(opts)
	zap.ReplaceGlobals(logger)
	zap.RedirectStdLog(logger)

	return func() {
		_ = logger.Sync()
	}
}
