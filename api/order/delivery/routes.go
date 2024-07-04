package customerRoutes

import (
	orderRoute "github.com/ahsansandiah/dpo-test/api/order/delivery/route"
	"github.com/ahsansandiah/dpo-test/packages/manager"
	"github.com/gorilla/mux"
)

func NewRoutes(r *mux.Router, mgr manager.Manager) {
	api := r.PathPrefix("").Subrouter()

	orderRoute.NewOrderRoute(mgr, api)
}
