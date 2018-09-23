package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	v1 := router.Group("/v1")
	{
		v1.GET("/creds", getCreds)
		v1.GET("/health", healthEndpoint)
	}

	srv := &http.Server{
		Addr:    ":8084",
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")

}

func getCreds(c *gin.Context) {

	credsFile := os.Getenv("CREDS_FILE_PATH")

	dat, err := ioutil.ReadFile(credsFile)
	if err != nil {

		c.JSON(503, gin.H{
			"status": "Creds file can't be read",
		})
	}

	c.String(http.StatusOK, string(dat))

}

func healthEndpoint(c *gin.Context) {

	c.JSON(200, gin.H{
		"status": "healthy",
	})
}