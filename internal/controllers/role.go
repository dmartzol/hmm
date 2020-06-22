package controllers

import (
	"log"
	"net/http"

	"github.com/dmartzol/hmmm/internal/models"
	"github.com/dmartzol/hmmm/pkg/httpresponse"
)

type roleStorage interface {
	RoleExists(name string) (bool, error)
	CreateRole(name string) (*models.Role, error)
}

func (api API) CreateRole(w http.ResponseWriter, r *http.Request) {
	var req models.CreateRoleReq
	err := httpresponse.Unmarshal(r, &req)
	if err != nil {
		log.Printf("JSON: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	exists, err := api.RoleExists(req.Name)
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if exists {
		log.Printf("Role '%s' already exists", req.Name)
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}
	role, err := api.storage.CreateRole(req.Name)
	if err != nil {
		log.Printf("error creating role: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	httpresponse.RespondJSON(w, role.View(nil))
}
