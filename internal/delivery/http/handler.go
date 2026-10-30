package httpdelivery

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/yohannes/sint-new-backend/internal/application"
	"github.com/yohannes/sint-new-backend/internal/domain"
)

type Handler struct {
	health       application.DatabaseHealth
	auth         application.Authenticator
	uploadPhoto  application.UploadPhoto
	syncUser     application.SyncUser
	listProducts application.ListProducts
	submitPrice  application.SubmitPrice
}

func NewHandler(
	health application.DatabaseHealth,
	auth application.Authenticator,
	uploadPhoto application.UploadPhoto,
	syncUser application.SyncUser,
	listProducts application.ListProducts,
	submitPrice application.SubmitPrice,
) http.Handler {
	handler := &Handler{
		health:       health,
		auth:         auth,
		uploadPhoto:  uploadPhoto,
		syncUser:     syncUser,
		listProducts: listProducts,
		submitPrice:  submitPrice,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.healthCheck)
	mux.Handle(
		"GET /api/v1/me",
		requireUser(auth, http.HandlerFunc(handler.currentUser)),
	)
	mux.Handle(
		"POST /api/v1/uploads/photos",
		requireUser(auth, http.HandlerFunc(handler.upload)),
	)
	mux.Handle(
		"GET /api/v1/products",
		requireUser(auth, http.HandlerFunc(handler.products)),
	)
	mux.Handle(
		"POST /api/v1/price-updates",
		requireUser(auth, http.HandlerFunc(handler.createPriceUpdate)),
	)
	return withCORS(mux)
}

func (handler *Handler) createPriceUpdate(writer http.ResponseWriter, request *http.Request) {
	user, ok := userFromContext(request.Context())
	if !ok {
		writeJSON(writer, http.StatusInternalServerError, map[string]string{
			"error": "authenticated user missing from context",
		})
		return
	}
	var update domain.PriceUpdate
	if err := json.NewDecoder(request.Body).Decode(&update); err != nil {
		writeJSON(writer, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}
	if err := handler.submitPrice.Execute(request.Context(), user, update); err != nil {
		if errors.Is(err, application.ErrInvalidPriceUpdate) {
			writeJSON(writer, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		log.Printf("create price update: %v", err)
		writeJSON(writer, http.StatusInternalServerError, map[string]string{
			"error": "could not submit price update",
		})
		return
	}
	writeJSON(writer, http.StatusCreated, map[string]string{"status": "pending"})
}

func (handler *Handler) healthCheck(writer http.ResponseWriter, request *http.Request) {
	if err := handler.health.Ping(request.Context()); err != nil {
		writeJSON(writer, http.StatusServiceUnavailable, map[string]string{
			"status": "database_unavailable",
		})
		return
	}
	writeJSON(writer, http.StatusOK, map[string]string{"status": "ok"})
}

func (handler *Handler) currentUser(writer http.ResponseWriter, request *http.Request) {
	user, ok := userFromContext(request.Context())
	if !ok {
		writeJSON(writer, http.StatusInternalServerError, map[string]string{
			"error": "authenticated user missing from context",
		})
		return
	}
	syncedUser, err := handler.syncUser.Execute(request.Context(), user)
	if err != nil {
		log.Printf("sync authenticated user: %v", err)
		writeJSON(writer, http.StatusInternalServerError, map[string]string{
			"error": "could not synchronize authenticated user",
		})
		return
	}
	writeJSON(writer, http.StatusOK, syncedUser)
}

func (handler *Handler) products(writer http.ResponseWriter, request *http.Request) {
	products, err := handler.listProducts.Execute(request.Context())
	if err != nil {
		log.Printf("list products: %v", err)
		writeJSON(writer, http.StatusInternalServerError, map[string]string{
			"error": "could not load products",
		})
		return
	}
	writeJSON(writer, http.StatusOK, products)
}

func (handler *Handler) upload(writer http.ResponseWriter, request *http.Request) {
	const maxUploadSize = 10 << 20
	request.Body = http.MaxBytesReader(writer, request.Body, maxUploadSize)
	if err := request.ParseMultipartForm(maxUploadSize); err != nil {
		writeJSON(writer, http.StatusBadRequest, map[string]string{
			"error": "photo must be a multipart file smaller than 10 MB",
		})
		return
	}
	file, header, err := request.FormFile("photo")
	if err != nil {
		writeJSON(writer, http.StatusBadRequest, map[string]string{
			"error": "multipart field 'photo' is required",
		})
		return
	}
	defer file.Close()

	url, err := handler.uploadPhoto.Execute(
		request.Context(),
		file,
		header.Filename,
		header.Header.Get("Content-Type"),
	)
	if err != nil {
		if errors.Is(err, application.ErrInvalidPhoto) {
			writeJSON(writer, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		log.Printf("upload photo: %v", err)
		writeJSON(writer, http.StatusBadGateway, map[string]string{
			"error": "photo storage is temporarily unavailable",
		})
		return
	}
	writeJSON(writer, http.StatusCreated, map[string]string{"url": url})
}
