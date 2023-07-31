package app

import (
	"github.com/darrylmorton/ct-iot-event-service/internal/data"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (app *application) healthCheck(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	result := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

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

	_, err := uuid.Parse(id)
	if err != nil {
		app.logger.Printf("invalid uuid: %s", id)
		http.Error(c.Writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	c.Header("Content-Type", "application/json")

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

func (app *application) putEvent(c *gin.Context) {
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

	c.Header("Content-Type", "application/json")

	var input = data.Event{}

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
