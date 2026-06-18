package main

import (
	"SOCIAL/internal/store"
	"fmt"
	"math/rand"
	"net/http"
)

//type commentKey string
//
//const commentCtx postKey = "comment"

type CreateCommentPayload struct {
	Content   string `json:"content" validate:"required,max=1000"`
	CreatedAt string `json:"created_at"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateCommentPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// check if payload content is empty
	if payload.Content == "" {
		app.badRequestResponse(w, r, fmt.Errorf("content is req"))
		return
	}

	// Get the PostID from the dB

	userMin := 1
	userMax := 100

	userID := rand.Intn(userMax-userMin+1) + userMin

	var finalUserID int64 = int64(userID)

	// Get the UserID from the dB

	postMin := 1
	postMax := 100

	postID := rand.Intn(postMax-postMin+1) + postMin

	var finalPostID int64 = int64(postID)

	comment := &store.Comment{
		Content: payload.Content,
		UserID:  finalUserID,
		// Change after auth
		PostID: finalPostID,
	}

	ctx := r.Context()

	if err := app.store.Comments.Create(ctx, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
