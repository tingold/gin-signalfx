package middleware

import (
	"testing"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	"context"
)

func TestDefaultConfig(t *testing.T) {

}

func TestHealthMonitor_Datapoints(t *testing.T) {

}

func TestSignalFx(t *testing.T) {

	g := gin.New()
	g.Use(SignalFx(Config{SignalFXKey:"INSERT KEY HERE"}))



	srv := &http.Server{
		Addr:    ":8000",
		Handler: g,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	time.Sleep(10 * time.Second)
	srv.Shutdown(context.Background())

}
