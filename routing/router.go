package routing

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"consumer/models"
)

type module interface {
	ServeBatch(batch models.Batch) error
	GetBufferFreeSpace() (int, error)
}

type Handler struct {
	m module
}

func NewHandler(m module) *Handler {
	return &Handler{m: m}
}

func (h Handler) GetBuffer(c echo.Context) error {
	bufferFreeSpace, err := h.m.GetBufferFreeSpace()
	if err != nil {
		return c.String(http.StatusServiceUnavailable, strconv.Itoa(bufferFreeSpace))
	}
	return c.String(http.StatusOK, strconv.Itoa(bufferFreeSpace))
}

func (h Handler) PostBatch(c echo.Context) error {
	batch, err := parseRequestBody(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	err = h.m.ServeBatch(batch)
	if err != nil {
		return c.String(http.StatusServiceUnavailable, err.Error())
	}

	return c.String(http.StatusOK, "Ok")
}

func parseRequestBody(c echo.Context) (models.Batch, error) {

	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = c.Request().Body.Close(); err != nil {
			log.Error(err.Error())
		}
	}()

	var batch models.Batch
	err = json.Unmarshal(b, &batch)
	if err != nil {
		return nil, err
	}

	return batch, err
}
