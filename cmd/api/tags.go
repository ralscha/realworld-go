package main

import (
	"net/http"
	"realworldgo.rasc.ch/cmd/api/dto"
	"realworldgo.rasc.ch/internal/models"
	"realworldgo.rasc.ch/internal/response"
)

func (app *application) tagsGet(w http.ResponseWriter, r *http.Request) {
	tags, err := models.Tags().All(r.Context(), app.db)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	var tagsResponse = dto.TagsWrapper{
		Tags: make([]string, len(tags)),
	}

	for i, tag := range tags {
		tagsResponse.Tags[i] = tag.Tag
	}

	response.JSON(w, http.StatusOK, tagsResponse)
}
