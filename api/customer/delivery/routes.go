package customerRoutes

import (
	customerRoute "github.com/ahsansandiah/dpo-test/api/customer/delivery/route"
	"github.com/ahsansandiah/dpo-test/packages/manager"
	"github.com/gorilla/mux"
)

func NewRoutes(r *mux.Router, mgr manager.Manager) {
	api := r.PathPrefix("").Subrouter()

	customerRoute.NewCustomerRoute(mgr, api)
}
