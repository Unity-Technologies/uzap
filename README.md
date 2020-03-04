# uzap
[![GoDoc](https://godoc.prd.cds.internal.unity3d.com/github.com/Unity-Technologies/uzap?status.svg)](https://godoc.prd.cds.internal.unity3d.com/github.com/Unity-Technologies/uzap)

This repo contains basic configuration for go.uber.org/zap, for easy use with 12-factor apps.

### Usage
``` go
import (
	"go.uber.org/zap"
	"github.com/Unity-Technologies/uzap"
)
zap.ReplaceGlobals(uzap.Log)
defer uzap.Log.Sync()
logger.Info("failed to fetch URL",
	// Structured context as strongly typed Field values.
	zap.String("url", url),
	zap.Int("attempt", 3),
	zap.Duration("backoff", time.Second),
)
```