# :zap: uzap [![GoDoc][doc-img]][doc]

This repo contains basic configuration for go.uber.org/zap, for easy use with 12-factor apps.

### Usage
``` go
import (
	"go.uber.org/zap"
	"github.com/Unity-Technologies/uzap"
)

// For ease of use, call this in main()
zap.ReplaceGlobals(uzap.Log)
defer uzap.Log.Sync()

// Example
zap.L().Info("failed to fetch URL",
	// Structured context as strongly typed Field values.
	zap.String("url", url),
	zap.Int("attempt", 3),
	zap.Duration("backoff", time.Second),
)
```

### Features
On it's own this does mostly nothing, [go.uber.org/zap](https://pkg.go.dev/go.uber.org/zap) is the logger, this just makes a few assumptions about how you'd use zap to make it even easier. 
- It uses a global variable to reference the logger, which allows it to be used from anywhere without prior knowledge, ie. just call `zap.L()` without worrying about the details.
- It uses a package init function to ensure that the logger is ready to use, even if you call it with-in the first few lines of your program code.
- It sends informational logs to stdout and higher priority logs to stderr, which is needed for production-ready kubernetes logging.
- It encodes the output in a stackdriver friendly (but not overly specific) (json) format.
- It has a global atomic error log level, so it can be safely modified from anywhere.
- It optionally supports a human friendly debug output.
- It reads config from environment variables (can be changed/overridden afterwards).

[doc-img]: https://img.shields.io/badge/godoc-reference-blue
[doc]: https://godoc.prd.cds.internal.unity3d.com/github.com/Unity-Technologies/uzap
