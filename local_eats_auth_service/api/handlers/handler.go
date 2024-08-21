package handlers

import (
	"auth_serice/storage"
)

type Handler struct {
	userStorage storage.IStorage
}

func (h *Handler) Users() storage.IUserStorage {
	return h.userStorage.Users()
}

func (h *Handler) Kitchens() storage.IKitchenStorage {
	return h.userStorage.Kitchens()
}

func NewHandler(userStorage storage.IStorage) *Handler {
	return &Handler{userStorage: userStorage}
}
