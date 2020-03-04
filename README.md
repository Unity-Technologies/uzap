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


[doc-img]: https://img.shields.io/badge/godoc-reference-blue
[doc]: https://godoc.prd.cds.internal.unity3d.com/github.com/Unity-Technologies/uzap
