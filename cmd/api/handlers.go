package main

import (
	"net/http"
	"realworldgo.rasc.ch/internal/response"
)

func (app *application) status(w http.ResponseWriter, _ *http.Request) {
	data := map[string]string{
		"Status": "OK",
	}

	response.JSON(w, http.StatusOK, data)
}
