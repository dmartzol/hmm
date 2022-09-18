package service

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type service struct {
	name       string
	apiHandler http.Handler
}

func NewService(name string, structuredLogging bool) *service {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	service := service{
		name:       name,
		apiHandler: r,
	}

	return &service
}
