package main

import (
	"github.com/Shankara130/compressor/internal/delivery/http"
	"github.com/Shankara130/compressor/internal/infrastructure/queue"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	r := gin.Default()

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	jobQueue := queue.NewRedisQueue(redisClient)
	handler := http.NewHandler(jobQueue)

	r.POST("/upload", handler.Upload)
	r.GET("/jobs/:id", handler.GetStatus)
	r.GET("/download/:id", handler.Download)

	r.Run(":8080")
}
