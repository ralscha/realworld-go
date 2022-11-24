package main

import (
	"database/sql"
	"net/http"
	"realworldgo.rasc.ch/cmd/api/dto"
	"realworldgo.rasc.ch/internal/models"
	"realworldgo.rasc.ch/internal/response"
)

func (app *application) tagsGet(w http.ResponseWriter, r *http.Request) {
	tx := r.Context().Value(transactionKey).(*sql.Tx)
	tags, err := models.Tags().All(r.Context(), tx)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	var tagsResponse = dto.TagsWrapper{
		Tags: make([]string, len(tags)),
	}

	for i, tag := range tags {
		tagsResponse.Tags[i] = tag.Name
	}

	response.JSON(w, http.StatusOK, tagsResponse)
}
