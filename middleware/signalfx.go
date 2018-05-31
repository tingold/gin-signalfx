package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/signalfx/golib/sfxclient"
	"github.com/signalfx/golib/datapoint"
	"golang.org/x/net/context"
	"time"
	"os"
	"fmt"
	"sync/atomic"
)

type Config struct {
	SignalFXKey string
	ServiceName string
}

type GinMonitor struct {
	id             map[string]string
	get            int64
	post           int64
	delete         int64
	patch          int64
	put            int64
	head           int64
	hits           int64
	responseErrors int64

	responseTime *sfxclient.RollingBucket
	responseSize *sfxclient.RollingBucket
	requestSize  *sfxclient.RollingBucket
}

func (c *GinMonitor) Datapoints() []*datapoint.Datapoint {

	dps := []*datapoint.Datapoint{sfxclient.Cumulative("gin.request.get", c.id, atomic.SwapInt64(&c.get, 0)),
		sfxclient.Cumulative("gin.request.head", c.id, atomic.SwapInt64(&c.head, 0)),
		sfxclient.Cumulative("gin.request.put", c.id, atomic.SwapInt64(&c.put, 0)),
		sfxclient.Cumulative("gin.request.post", c.id, atomic.SwapInt64(&c.post, 0)),
		sfxclient.Cumulative("gin.request.delete", c.id, atomic.SwapInt64(&c.delete, 0)),
		sfxclient.Cumulative("gin.request.hits", c.id, atomic.SwapInt64(&c.hits, 0)),
		sfxclient.Cumulative("gin.request.errors", c.id, atomic.SwapInt64(&c.responseErrors, 0))}

	mc := sfxclient.NewMultiCollector(c.requestSize, c.responseSize, c.responseTime)

	dps = append(dps, sfxclient.GoMetricsSource.Datapoints()...)
	return append(dps, mc.Datapoints()...)
}

func DefaultConfig() Config {
	return Config{
		SignalFXKey: "invalid",
		ServiceName: "unknown",
	}
}

func SignalFx(config Config) gin.HandlerFunc {

	hostname, _ := os.Hostname()
	id := make(map[string]string)
	id["host.name"] = hostname
	id["service.name"] = config.ServiceName

	gm := GinMonitor{
		id:             id,
		responseTime:   sfxclient.NewRollingBucket("gin.response.time", id),
		responseSize:   sfxclient.NewRollingBucket("gin.response.size", id),
		requestSize:    sfxclient.NewRollingBucket("gin.request.size", id),
		get:            int64(0),
		post:           int64(0),
		delete:         int64(0),
		patch:          int64(0),
		put:            int64(0),
		head:           int64(0),
		hits:           int64(0),
		responseErrors: int64(0),
	}

	scheduler := sfxclient.NewScheduler()
	sink := sfxclient.NewHTTPSink()
	sink.AuthToken = config.SignalFXKey
	scheduler.Sink = sink
	scheduler.AddCallback(&gm)

	go scheduler.Schedule(context.Background())

	return func(c *gin.Context) {

		fmt.Println("setting values")
		start := time.Now()
		reqSz := c.Request.ContentLength
		verb := c.Request.Method
		//continue processing
		c.Next()
		//operating on response now
		dur := time.Since(start)

		stat := c.Writer.Status()
		if stat > 399 {
			atomic.AddInt64(&gm.responseErrors, 1)
		}
		respSz := c.Writer.Size()

		gm.responseTime.Add(dur.Seconds())
		gm.responseSize.Add(float64(respSz))
		gm.requestSize.Add(float64(reqSz))

		atomic.AddInt64(&gm.hits, 1)

		fmt.Println("verb:" + verb)

		switch verb {
		case "POST":
			atomic.AddInt64(&gm.post, 1)
			break
		case "GET":
			atomic.AddInt64(&gm.get, 1)
			break
		case "PUT":
			atomic.AddInt64(&gm.put, 1)
			break
		case "PATCH":
			atomic.AddInt64(&gm.patch, 1)
			break
		case "DELETE":
			atomic.AddInt64(&gm.delete, 1)
			break
		case "HEAD":
			atomic.AddInt64(&gm.head, 1)
			break
		}

	}

}
