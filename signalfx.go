package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/signalfx/golib/sfxclient"
	"github.com/signalfx/golib/datapoint"
	"golang.org/x/net/context"
)



type Config struct {
	SignalFXKey string
}

type HealthMonitor struct {


}



func (c *HealthMonitor) Datapoints() []*datapoint.Datapoint {
	return sfxclient.GoMetricsSource.Datapoints()
}

func DefaultConfig() Config {
	return Config{
		SignalFXKey: "invalid",
	}
}

func SignalFx(config Config) gin.HandlerFunc {

	scheduler := sfxclient.NewScheduler()
	sink := sfxclient.NewHTTPSink()
	sink.AuthToken = config.SignalFXKey
	scheduler.Sink = sink
	scheduler.AddGroupedCallback("go-stats",sfxclient.GoMetricsSource)
	go scheduler.Schedule(context.Background())

	return func(c *gin.Context) {



	}

}
