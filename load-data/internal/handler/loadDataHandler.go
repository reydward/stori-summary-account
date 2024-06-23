package handler

import (
	"fmt"
	"load-data/internal/repository"
	"net/http"
)

type LoadDataHandler struct {
	repo repository.LoadDataRepository
}

func NewLoadDataHandler(repo repository.LoadDataRepository) *LoadDataHandler {
	return &LoadDataHandler{repo: repo}
}

func (h *LoadDataHandler) Health(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<h1>The Lambda Load Data is working!<h1>\n")
}
