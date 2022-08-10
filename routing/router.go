package routing

import (
	"consumer/models"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"log"
	"net/http"
)

type module interface {
	ServeBatch(batch models.Batch) error
}

type Handler struct {
	m module
}

func NewHandler(m module) *Handler {
	return &Handler{m: m}
}

func (h Handler) PostBatch(c echo.Context) error {
	batch, err := parseRequestBody(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	err = h.m.ServeBatch(batch)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "Ok")
}

func parseRequestBody(c echo.Context) (models.Batch, error) {
	var batch models.Batch

	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = c.Request().Body.Close(); err != nil {
			log.Printf(err.Error())
		}
	}()

	err = json.Unmarshal(b, &batch)
	if err != nil {
		return nil, err
	}

	return batch, err
}
