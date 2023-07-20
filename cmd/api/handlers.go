package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (app *application) healthCheck(c *gin.Context) {
	result := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	c.Header("Content-Type", "application/json")

	c.IndentedJSON(200, result)
}

func (app *application) getEvents(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	results, err := app.models.Events.GetEvents()
	if err != nil {
		http.Error(c.Writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(200, results)
}

func (app *application) getEvent(c *gin.Context) {
	id := c.Param("id")

	app.logger.Println("ID", id)

	c.Header("Content-Type", "application/json")

	// TODO replace status codes with error types
	results, err, statusCode := app.models.Events.GetEvent(id)
	if err != nil {
		switch {
		case statusCode == 400:
			http.Error(c.Writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		case statusCode == 404:
			http.Error(c.Writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		default:
			http.Error(c.Writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	c.IndentedJSON(200, results)
}
