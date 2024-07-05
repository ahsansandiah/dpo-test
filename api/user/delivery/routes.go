package userRoutes

import (
	userRoute "github.com/ahsansandiah/dpo-test/api/user/delivery/route"
	"github.com/ahsansandiah/dpo-test/packages/manager"
	"github.com/gorilla/mux"
)

func NewRoutes(r *mux.Router, mgr manager.Manager) {
	api := r.PathPrefix("").Subrouter()

	userRoute.NewUserRoute(mgr, api)
}
