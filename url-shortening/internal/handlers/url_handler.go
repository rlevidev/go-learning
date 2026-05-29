package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	model "url-shortening/internal/models"
	"url-shortening/internal/service"
)

type UrlHandler struct {
	service *service.UrlService
}

func New(service *service.UrlService) *UrlHandler {
	return &UrlHandler{
		service: service,
	}
}

func (h *UrlHandler) PostShorten(w http.ResponseWriter, r *http.Request) {
	var url model.Url
	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		http.Error(w, "Requisição inválida", http.StatusBadRequest)
		return
	}

	if url.URLOriginal == "" {
		http.Error(w, "A URL original é obrigatória", http.StatusBadRequest)
		return
	}

	createdUrl, err := h.service.CreateShortURL(r.Context(), url.URLOriginal)
	if err != nil {
		log.Printf("Erro ao criar uma url encurtada: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUrl)
}

func (h *UrlHandler) GetShorten(w http.ResponseWriter, r *http.Request) {
	urls, err := h.service.GetShortURL(r.Context(), r.PathValue("shortCode"))
	if err != nil {
		// Não achou no banco, return 404
		if err == sql.ErrNoRows {
			http.Error(w, "URL não encontrada", http.StatusNotFound)
			return
		}

		// Qualquer outro erro (ex: banco off) return 500
		http.Error(w, "Erro ao buscar a URL original", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, urls.UrlOriginal, http.StatusFound)
}

func (h *UrlHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	urls, err := h.service.GetShortURL(r.Context(), r.PathValue("shortCode"))
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "URL não encontrada", http.StatusNotFound)
			return
		}

		http.Error(w, "Erro ao buscar a URL original", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(urls)
}

func (h *UrlHandler) DeleteUrlShorted(w http.ResponseWriter, r *http.Request) {
	if err := h.service.DeleteShortURL(r.Context(), r.PathValue("shortCode")); err != nil {
		http.Error(w, "Erro ao excluir a URL encurtada", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UrlHandler) UpdateShortUrl(w http.ResponseWriter, r *http.Request) {
	var url model.Url
	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		http.Error(w, "Requisição inválida", http.StatusBadRequest)
		return
	}
	
	updatedUrl, err := h.service.UpdateShortURL(r.Context(), url.ShortCode, url.URLOriginal)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "URL não encontrada", http.StatusNotFound)
			return
		}
		
		log.Printf("Erro ao modificar uma URL encurtada: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedUrl)
}
