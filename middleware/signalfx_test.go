package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestSignalFx(t *testing.T) {

	g := gin.New()
	g.Use(SignalFx(Config{SignalFXKey: "INSERT KEY HERE", ServiceName: "test-code"}))

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
