package post

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func CreatePostHandler(service *PostService, w http.ResponseWriter, r *http.Request) {
	var request CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Error en el formato del cuerpo", http.StatusBadRequest)
		return
	}
	err := service.CreatePost(request)
	if err != nil {
		http.Error(w, "Error al crear el post", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetPostByIDHandler(service *PostService, w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	post, err := service.FindPostById(id)
	if err != nil {
		http.Error(w, "Error al buscar el post", http.StatusNotFound)
		return
	}
	if post.Id == uuid.Nil {
		http.Error(w, "Post no encontrado", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}
