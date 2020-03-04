package uzap

import (
	"os"

	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/kelseyhightower/envconfig"
)

var (
	// Log is global logger
	Log   *zap.Logger
	Level zap.AtomicLevel
)

type Config struct {
	Level zapcore.Level `required:"true" default:"warn"`
	Debug bool          `required:"true" default:"false"`
}

// Use package init to avoid race conditions for GRPC options
// sync.Once still suffers from races, init functions are less complex than sync.once + waitgroup
func init() {
	var cfg Config
	if err := envconfig.Process("log", &cfg); err != nil {
		panic(err)
	}

	Level = zap.NewAtomicLevelAt(cfg.Level)

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	// It is useful for Kubernetes deployment.
	// Kubernetes interprets os.Stdout log items as INFO and os.Stderr log items
	// as ERROR by default.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= Level.Level() && lvl < zapcore.ErrorLevel
	})

	// Output channels
	consoleInfos := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	// Setup Config
	var (
		ecfg zapcore.EncoderConfig
		enc  zapcore.Encoder
	)

	if cfg.Debug {
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
	Log = zap.New(core)
}
