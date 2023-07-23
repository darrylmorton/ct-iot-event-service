package main

import (
	"github.com/gin-gonic/gin"
)

func (app *application) router() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")

	v1.GET("/health", app.healthCheck)
	v1.GET("/events", app.getEvents)
	v1.GET("/events/:id", app.getEvent)
	v1.PUT("/events/:id", app.putEvent)

	return r
}
