package main

type Handler struct {
	store *LocationStore
}

func NewHandler(store *LocationStore) *Handler {
	return &Handler{store: store}
}

