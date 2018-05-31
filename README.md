# gin-signalfx
SignalFx Middleware for Gin

[![Go Report Card](https://goreportcard.com/badge/github.com/tingold/gin-signalfx)](https://goreportcard.com/report/github.com/tingold/gin-signalfx)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Super simple usage: 
```	
import ( github.com/tingold/gin-signalfx/middleware
         github.com/gin-gonic/gin
        )
g := gin.New()
g.Use(SignalFx(Config{SignalFXKey:"INSERT KEY HERE", ServiceName: "YOUR GIN SERVICE"}))

//start gin normally
...
```
Upon starting the server you should start to see metrics appear in your SignalFx dashboard

## Metrics 

All Gin specific metrics can be found by searching for `gin.*` in SignalFx

The following metrics will be reported:

- Counts for each of the HTTP verbs in the format of `gin.request.{verb}`
- Hits (e.g requests)
- Errors: response codes > 399
- Response times
- Request / Response Sizes

The following Go runtime metrics are populated (which should be self explanatory):
``` 
Alloc
BuckHashSys
Frees
GCSys
GOMAXPROCS
HeapAlloc
HeapIdle
HeapInuse
HeapObjects
HeapReleased
HeapSys
LastGC
Lookups
Mallocs
MCacheInuse
MCacheSys
MSpanInuse
MSpanSys
NextGC
num_cgo_call
num_cpu
num_goroutine
NumGC
OtherSys
PauseTotalNs
process.uptime.ns
StackInuse
StackSys
Sys
TotalAlloc 
```