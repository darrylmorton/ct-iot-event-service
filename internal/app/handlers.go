package app

import (
	"github.com/darrylmorton/ct-iot-event-service/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (app *Application) healthCheck(c *gin.Context) {
	c.Header("Content-Type", "Application/json")

	result := map[string]string{
		"status":      "available",
		"environment": app.config.Env,
		"version":     Version,
	}

	c.IndentedJSON(200, result)
}

func (app *Application) getEvents(c *gin.Context) {
	c.Header("Content-Type", "Application/json")

	results, err := app.models.Events.GetEvents()
	if err != nil {
		http.Error(c.Writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(200, results)
}

func (app *Application) getEvent(c *gin.Context) {
	id := c.Param("id")

	_, err := uuid.Parse(id)
	if err != nil {
		app.logger.Printf("invalid uuid: %s", id)
		http.Error(c.Writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	c.Header("Content-Type", "Application/json")

	results, err, statusCode := app.models.Events.GetEvent(id)
	if err != nil {
		switch {
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

func (app *Application) putEvent(c *gin.Context) {
	id := c.Param("id")

	_, err := uuid.Parse(id)
	if err != nil {
		app.logger.Printf("invalid uuid: %s", id)
		http.Error(c.Writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	_, getEventErr, getEventStatusCode := app.models.Events.GetEvent(id)
	if getEventErr != nil {
		http.Error(c.Writer, http.StatusText(getEventStatusCode), getEventStatusCode)
		return
	}

	c.Header("Content-Type", "Application/json")

	var input = models.Event{}

	err = app.readJSON(c.Writer, c.Request, &input)
	if err != nil {
		http.Error(c.Writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	results, err, statusCode := app.models.Events.PutEvent(id, input)
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
